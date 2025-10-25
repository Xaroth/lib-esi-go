package ratelimiter

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidHeader = errors.New("invalid header")
)

type Group struct {
	Name       string
	BucketSize int
	WindowSize time.Duration
}

func NewGroup(name string, header string) (*Group, error) {
	if len(name) == 0 || len(header) == 0 {
		return nil, ErrInvalidHeader
	}

	parts := strings.Split(header, "/")
	if len(parts) != 2 {
		return nil, ErrInvalidHeader
	}
	bucketSize, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, ErrInvalidHeader
	}
	windowSize, err := time.ParseDuration(parts[1])
	if err != nil {
		return nil, ErrInvalidHeader
	}
	return &Group{
		Name:       name,
		BucketSize: bucketSize,
		WindowSize: windowSize,
	}, nil
}

func (g *Group) TimePerToken() time.Duration {
	milliseconds := g.WindowSize.Milliseconds() / int64(g.BucketSize)
	return time.Duration(milliseconds) * time.Millisecond
}

func (g *Group) TargetSize(percentage float64) int {
	return int(math.Floor(float64(g.BucketSize) * percentage))
}
