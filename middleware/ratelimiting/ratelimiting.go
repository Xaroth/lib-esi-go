package ratelimiting

import (
	"context"
	"errors"
	"net/http"

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
	ErrInvalidHeader = errors.New("invalid header")
)

type RateLimiter interface {
	Schedule(ctx context.Context, info request.RequestInfo, forOwner int64) (func(*http.Response), error)
}
