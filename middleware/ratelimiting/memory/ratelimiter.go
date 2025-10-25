package memory

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/xaroth/lib-esi-go/middleware/ratelimiting"
	"github.com/xaroth/lib-esi-go/request"
)

var (
	maximumWaitTimePerAttempt = 10 * time.Second
)

type memoryRateLimiter struct {
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

func New() ratelimiting.RateLimiter {
	return &memoryRateLimiter{
		targetPercentage:          ratelimiting.DefaultRateLimitTargetPercentage,
		estimatedTokensPerRequest: ratelimiting.DefaultEstimatedTokensPerRequest,

		groups:     make(map[string]*Group),
		patternMap: make(map[string]*Group),
		buckets:    make(map[bucketKey]*Bucket),
	}
}

func (r *memoryRateLimiter) getPathGroup(info request.RequestInfo) *Group {
	group, _ := r.patternMap[info.Path]
	return group
}

func (r *memoryRateLimiter) setPathGroup(info request.RequestInfo, group *Group) {
	r.patternMu.Lock()
	defer r.patternMu.Unlock()
	r.patternMap[info.Path] = group
}

func (r *memoryRateLimiter) updateGroup(group *Group) *Group {
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

func (r *memoryRateLimiter) ListBuckets() []*BucketInfo {
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

func (r *memoryRateLimiter) getBucket(group *Group, owner int64) *Bucket {
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

func (r *memoryRateLimiter) CleanupExpiredBuckets() {
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

func (r *memoryRateLimiter) delayRequest(ctx context.Context, bucket *Bucket) {
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

func (r *memoryRateLimiter) processResponse(info request.RequestInfo, owner int64, resp *http.Response) {
	currentGroup := r.getPathGroup(info)

	if resp == nil {
		return
	}

	groupHeader := resp.Header.Get(ratelimiting.RateLimitGroupKey)
	limitHeader := resp.Header.Get(ratelimiting.RateLimitLimitKey)
	if groupHeader == "" || limitHeader == "" {
		return
	}

	group, err := NewGroup(groupHeader, limitHeader)
	if err != nil {
		return
	}
	group = r.updateGroup(group)

	if currentGroup == nil {
		r.setPathGroup(info, group)
	} else if group.Name != currentGroup.Name {
		r.setPathGroup(info, group)
	}

	remainingHeader := resp.Header.Get(ratelimiting.RateLimitRemainingKey)
	if remainingHeader == "" {
		return
	}

	remaining, err := strconv.Atoi(remainingHeader)
	if err != nil {
		return
	}

	bucket := r.getBucket(group, owner)
	bucket.SetRemainingTokens(remaining)
}

// Schedule handles rate limiting for a request. It extracts token information
// and delays the request if necessary to respect rate limits.
// Returns a callback that should be called with the response to update rate limit state.
func (r *memoryRateLimiter) Schedule(ctx context.Context, info request.RequestInfo, owner int64) (func(*http.Response), error) {
	group := r.getPathGroup(info)
	var bucketOwner int64 = -1

	switch {
	// No authentication token provided. Request is shared across all applications from the same IP.
	case owner == 0:
		bucketOwner = -1
	// No scope required, but token provided. Request is bucketed to the application.
	case info.RequiredScope == "":
		bucketOwner = -2
	// Request is bucketed to the application and owner.
	default:
		bucketOwner = owner
	}

	bucket := r.getBucket(group, bucketOwner)

	if bucket != nil {
		if bucket.CurrentUsage() >= r.targetPercentage {
			// If we are over the target percentage, delay the request until we are under the target percentage.
			r.delayRequest(ctx, bucket)
		}
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return func(resp *http.Response) {
		r.processResponse(info, owner, resp)
	}, nil
}
