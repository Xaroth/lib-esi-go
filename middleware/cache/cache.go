package cache

import "time"

// Cache defines the interface for cache storage implementations.
type Cache interface {
	Get(key string) (*CacheEntry, bool)
	Set(key string, entry *CacheEntry)
	Delete(key string)
	Clear()
}

// Entry represents a cached HTTP response.
type CacheEntry struct {
	ETag         string
	LastModified time.Time
}
