package ratelimiting

import (
	"errors"
	"net/http"
	"time"
)

const (
	RateLimitGroupKey     = "X-Ratelimit-Group"
	RateLimitLimitKey     = "X-Ratelimit-Limit"
	RateLimitRemainingKey = "X-Ratelimit-Remaining"

	DefaultRateLimitTargetPercentage = 0.75
	DefaultEstimatedTokensPerRequest = 5
)

var (
	ErrInvalidHeader = errors.New("invalid header")
)

type BucketStatistics struct {
	// The name of the rate limit group
	Group string

	// The owner of the bucket
	// Either a character ID, or -1 for shared buckets, or -2 for application buckets.
	Owner int64

	// The number of tokens currently available in the bucket
	EffectiveTokens int

	// The time the last request was made
	LastRequest time.Time
}

type RateLimiter interface {
	// Schedule a request to be delayed until the rate limit is no longer exceeded.
	Schedule(req *http.Request) (func(*http.Response), error)

	// List all active buckets and their statistics.
	ListBuckets() []*BucketStatistics
}
