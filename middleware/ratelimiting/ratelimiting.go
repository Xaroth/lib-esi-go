package ratelimiting

import (
	"errors"
	"net/http"
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

type RateLimiter interface {
	Schedule(req *http.Request) (func(*http.Response), error)
}
