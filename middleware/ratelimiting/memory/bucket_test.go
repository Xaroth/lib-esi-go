package memory_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/xaroth/lib-esi-go/middleware/ratelimiting/memory"
)

func TestEffectiveTokens(t *testing.T) {
	t.Parallel()

	group := &memory.Group{
		Name:       "test",
		BucketSize: 100,
		WindowSize: 10 * time.Second,
	}

	testCases := []struct {
		name                 string
		currentTokens        int
		timeSinceLastRequest time.Duration
		expectedTokens       int
		expectedUsage        float64
	}{
		{
			name:                 "no time since last request",
			currentTokens:        0,
			timeSinceLastRequest: 0,
			expectedTokens:       0,
			expectedUsage:        0.0,
		},
		{
			name:                 "10 seconds since last request",
			currentTokens:        100,
			timeSinceLastRequest: 10 * time.Second,
			expectedTokens:       0,
			expectedUsage:        0.0,
		},
		{
			name:                 "5 seconds since last request",
			currentTokens:        100,
			timeSinceLastRequest: 5 * time.Second,
			expectedTokens:       50,
			expectedUsage:        0.5,
		},
		{
			name:                 "3 seconds since last request",
			currentTokens:        100,
			timeSinceLastRequest: 3 * time.Second,
			expectedTokens:       70,
			expectedUsage:        0.7,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bucket := memory.FakeBucket(group, 1, tc.currentTokens, time.Now().Add(-tc.timeSinceLastRequest))
			effectiveTokens := bucket.EffectiveTokens()
			if diff := cmp.Diff(tc.expectedTokens, effectiveTokens); diff != "" {
				t.Fatalf("effective tokens mismatch (-want +got): %s", diff)
			}
			usage := bucket.CurrentUsage()
			if diff := cmp.Diff(tc.expectedUsage, usage); diff != "" {
				t.Fatalf("usage mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestTimeUntil(t *testing.T) {
	t.Parallel()

	group := &memory.Group{
		Name:       "test",
		BucketSize: 100,
		WindowSize: 10 * time.Second,
	}

	testCases := []struct {
		name                 string
		currentTokens        int
		timeSinceLastRequest time.Duration
		untilTokens          int
		expectedTime         time.Duration
	}{
		{
			name:          "immediately available from empty bucket",
			currentTokens: 0,
			untilTokens:   50,
			expectedTime:  0,
		},
		{
			name:                 "immediate available due to time since last request",
			currentTokens:        100,
			timeSinceLastRequest: 5 * time.Second,
			untilTokens:          50,
			expectedTime:         0,
		},
		{
			name:          "5 seconds until tokens",
			currentTokens: 100,
			untilTokens:   50,
			expectedTime:  5 * time.Second,
		},
		{
			name:          "10 seconds until tokens",
			currentTokens: 200,
			untilTokens:   100,
			expectedTime:  10 * time.Second,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bucket := memory.FakeBucket(group, 1, tc.currentTokens, time.Now().Add(-tc.timeSinceLastRequest))
			time := bucket.TimeUntil(tc.untilTokens)

			if diff := cmp.Diff(tc.expectedTime, time); diff != "" {
				t.Fatalf("time mismatch (-want +got): %s", diff)
			}
		})
	}
}
