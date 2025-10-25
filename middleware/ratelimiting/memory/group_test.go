package memory_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/xaroth/lib-esi-go/middleware/ratelimiting"
	"github.com/xaroth/lib-esi-go/middleware/ratelimiting/memory"
)

func TestNewGroup(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		nameHeader    string
		limitHeader   string
		expectedError error
		expectedGroup *memory.Group
	}{
		{
			name:          "valid header",
			nameHeader:    "test",
			limitHeader:   "100/10s",
			expectedError: nil,
			expectedGroup: &memory.Group{
				Name:       "test",
				BucketSize: 100,
				WindowSize: 10 * time.Second,
			},
		},
		{
			name:          "invalid header",
			nameHeader:    "test",
			limitHeader:   "100/10s/10s",
			expectedError: ratelimiting.ErrInvalidHeader,
			expectedGroup: nil,
		},
		{
			name:          "invalid bucket size",
			nameHeader:    "test",
			limitHeader:   "invalid/10s",
			expectedError: ratelimiting.ErrInvalidHeader,
			expectedGroup: nil,
		},
		{
			name:          "invalid window size",
			nameHeader:    "test",
			limitHeader:   "100/invalid",
			expectedError: ratelimiting.ErrInvalidHeader,
			expectedGroup: nil,
		},
		{
			name:          "invalid name",
			nameHeader:    "",
			limitHeader:   "100/10s",
			expectedError: ratelimiting.ErrInvalidHeader,
			expectedGroup: nil,
		},
		{
			name:          "invalid limit header",
			nameHeader:    "test",
			limitHeader:   "",
			expectedError: ratelimiting.ErrInvalidHeader,
			expectedGroup: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			group, err := memory.NewGroup(tc.nameHeader, tc.limitHeader)
			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, tc.expectedError) {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.expectedGroup, group); diff != "" {
				t.Fatalf("group mismatch (-want +got): %s", diff)
			}
		})
	}

}

func TestTimePerToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		group        *memory.Group
		expectedTime time.Duration
	}{
		{
			name: "valid group",
			group: &memory.Group{
				Name:       "test",
				BucketSize: 100,
				WindowSize: 10 * time.Second,
			},
			expectedTime: 100 * time.Millisecond,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			time := tc.group.TimePerToken()
			if diff := cmp.Diff(tc.expectedTime, time); diff != "" {
				t.Fatalf("time mismatch (-want +got): %s", diff)
			}
		})
	}
}

func TestTargetSize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		group        *memory.Group
		percentage   float64
		expectedSize int
	}{
		{
			name: "valid group",
			group: &memory.Group{
				Name:       "test",
				BucketSize: 100,
				WindowSize: 10 * time.Second,
			},
			percentage:   0.5,
			expectedSize: 50,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			size := tc.group.TargetSize(tc.percentage)
			if diff := cmp.Diff(tc.expectedSize, size); diff != "" {
				t.Fatalf("size mismatch (-want +got): %s", diff)
			}
		})
	}
}
