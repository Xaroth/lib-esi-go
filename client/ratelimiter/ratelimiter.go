package ratelimiter

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/xaroth/lib-esi-go/request"
)

const (
	RateLimitGroupKey     = "X-Ratelimit-Group"
	RateLimitLimitKey     = "X-Ratelimit-Limit"
	RateLimitRemainingKey = "X-Ratelimit-Remaining"

	DefaultRateLimitTargetPercentage = 0.75
	DefaultEstimatedTokensPerRequest = 5
)

var (
	maximumWaitTimePerAttempt = 10 * time.Second
)

type Ratelimiter struct {
	targetPercentage          float64
	estimatedTokensPerRequest int

	groups  map[string]*Group
	groupMu sync.Mutex

	patternMap map[string]*Group
	patternMu  sync.Mutex

	buckets   map[bucketKey]*Bucket
	bucketsMu sync.RWMutex
}

type BucketInfo struct {
	group           *Group
	owner           int64
	effectiveTokens int
	lastRequest     time.Time
}

var _ request.RequestScheduler = (*Ratelimiter)(nil)

func New() *Ratelimiter {
	return &Ratelimiter{
		// Aim for a percentage of the bucket size to be used before rate limiting.
		// This allows for burst requests, but still remains an overall good neighbour.
		targetPercentage: DefaultRateLimitTargetPercentage,
		// Assume that requests will always consume some tokens. This is to ensure
		// that we don't burst requests before we have been able to see how many tokens
		// we have remaining from responses.
		estimatedTokensPerRequest: DefaultEstimatedTokensPerRequest,
	}
}

func (r *Ratelimiter) getPathGroup(info request.RequestInfo) *Group {
	group, _ := r.patternMap[info.Path]
	return group
}

func (r *Ratelimiter) setPathGroup(info request.RequestInfo, group *Group) {
	r.patternMu.Lock()
	defer r.patternMu.Unlock()
	r.patternMap[info.Path] = group
}

func (r *Ratelimiter) updateGroup(group *Group) *Group {
	r.groupMu.Lock()
	defer r.groupMu.Unlock()

	if found, ok := r.groups[group.Name]; ok {
		found.BucketSize = group.BucketSize
		found.WindowSize = group.WindowSize
		return found
	}

	r.groups[group.Name] = group
	return group
}

func (r *Ratelimiter) ListBuckets() []*BucketInfo {
	info := make([]*BucketInfo, 0)
	r.bucketsMu.RLock()
	defer r.bucketsMu.RUnlock()
	for key, bucket := range r.buckets {
		info = append(info, &BucketInfo{
			group:           key.group,
			owner:           key.owner,
			effectiveTokens: bucket.EffectiveTokens(),
			lastRequest:     bucket.lastRequest,
		})
	}
	return info
}

func (r *Ratelimiter) getBucket(group *Group, owner int64) *Bucket {
	if group == nil {
		return nil
	}
	key := bucketKey{group: group, owner: owner}

	r.bucketsMu.RLock()
	bucket, ok := r.buckets[key]
	r.bucketsMu.RUnlock()

	if ok {
		return bucket
	}

	r.bucketsMu.Lock()
	defer r.bucketsMu.Unlock()
	bucket = NewBucket(group, owner)
	r.buckets[key] = bucket

	return bucket
}

func (r *Ratelimiter) CleanupExpiredBuckets() {
	keys := make([]bucketKey, 0)
	now := time.Now()

	r.bucketsMu.RLock()
	for key, bucket := range r.buckets {
		// If the bucket has not been used in the last two windows, we can safely assume it is back to 0 tokens.
		if now.Sub(bucket.lastRequest) > (bucket.group.WindowSize * 2) {
			keys = append(keys, key)
		}
	}
	r.bucketsMu.RUnlock()

	r.bucketsMu.Lock()
	defer r.bucketsMu.Unlock()
	for _, key := range keys {
		delete(r.buckets, key)
	}
}

func (r *Ratelimiter) delayRequest(ctx context.Context, bucket *Bucket) {
	target := bucket.group.TargetSize(r.targetPercentage)

	for {
		timeToWait := bucket.TimeUntil(target)
		if timeToWait <= 0 {
			// Done waiting.
			break
		}

		if timeToWait > maximumWaitTimePerAttempt {
			// Cap the wait time at the maximum time a request would take.
			// This ensures that very low rate limit groups don't cause
			// requests to wait for an unreasonable amount of time.
			// By the time the request ends we will have an accurate view of the
			// remaining tokens, and we can adjust the wait time accordingly.
			timeToWait = maximumWaitTimePerAttempt
		}

		select {
		case <-ctx.Done():
			// Context was cancelled, stop waiting.
			return
		case <-time.After(timeToWait):
			// Try again to ensure another request hasn't been made in the meanwhile.
			continue
		}
	}

	// Eagerly claim tokens for this request.
	bucket.Claim(r.estimatedTokensPerRequest)
}

func (r *Ratelimiter) processResponse(info request.RequestInfo, owner int64, resp *http.Response) {
	currentGroup := r.getPathGroup(info)

	if resp == nil {
		return
	}

	group, err := NewGroup(
		resp.Header.Get(RateLimitGroupKey),
		resp.Header.Get(RateLimitLimitKey),
	)
	if err != nil {
		return
	}
	group = r.updateGroup(group)

	if currentGroup == nil {
		r.setPathGroup(info, group)
	} else if group.Name != currentGroup.Name {
		r.setPathGroup(info, group)
	}

	remainingHeader := resp.Header.Get(RateLimitRemainingKey)
	remaining, err := strconv.Atoi(remainingHeader)
	if err != nil {
		return
	}

	bucket := r.getBucket(group, owner)
	bucket.Set(remaining)
}

func (r *Ratelimiter) Schedule(ctx context.Context, info request.RequestInfo, token request.Token) (func(*http.Response), error) {
	group := r.getPathGroup(info)
	var owner int64 = -1

	if token != nil {
		// If a token is sent, we will be placed on a separate bucket from other applications from the same IP.
		owner = -2

		// If a scope is required, the bucket key is dependent on the token owner.
		if info.RequiredScope != "" {
			owner = token.Owner()
		}
	}
	bucket := r.getBucket(group, owner)

	if bucket != nil {
		if bucket.CurrentUsage() >= r.targetPercentage {
			// If we are over the target percentage, delay the request until we are under the target percentage.
			r.delayRequest(ctx, bucket)
		}
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return func(resp *http.Response) {
		r.processResponse(info, owner, resp)
	}, nil
}
