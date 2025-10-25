package memory

import (
	"math"
	"sync"
	"time"
)

type Bucket struct {
	group         *Group
	currentTokens int
	lastRequest   time.Time
	mu            sync.RWMutex
}

func NewBucket(group *Group, owner int64) *Bucket {
	return &Bucket{
		group:         group,
		currentTokens: 0,
		lastRequest:   time.Now(),
	}
}

func FakeBucket(group *Group, owner int64, currentTokens int, lastRequest time.Time) *Bucket {
	return &Bucket{
		group:         group,
		currentTokens: currentTokens,
		lastRequest:   lastRequest,
	}
}

type bucketKey struct {
	group *Group
	owner int64
}

// Calculate the amoount of tokens we should currently be on by accounting for the time since last request.
func (b *Bucket) EffectiveTokens() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRequest).Seconds()
	windowsElapsed := elapsed / b.group.WindowSize.Seconds()
	tokensElapsed := float64(b.group.BucketSize) * windowsElapsed

	return int(math.Max(0, float64(b.currentTokens)-math.Floor(tokensElapsed)))
}

func (b *Bucket) CurrentUsage() float64 {
	return float64(b.EffectiveTokens()) / float64(b.group.BucketSize)
}

func (b *Bucket) TimeUntil(tokens int) time.Duration {
	missing := b.EffectiveTokens() - tokens
	if missing <= 0 {
		return 0
	}
	return b.group.TimePerToken() * time.Duration(missing)
}

func (b *Bucket) SetRemainingTokens(tokens int) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Calculate how many tokens have been used in the window based on the remaining tokens.
	b.currentTokens = b.group.BucketSize - tokens
	b.lastRequest = time.Now()

	return b.currentTokens
}

func (b *Bucket) Claim(tokens int) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.currentTokens += tokens
	b.lastRequest = time.Now()

	return b.currentTokens
}
