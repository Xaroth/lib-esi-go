package cache

import (
	"container/list"
	"sync"
)

// memoryStorage is an in-memory cache implementation with LRU eviction.
type memoryStorage struct {
	mu      sync.Mutex
	cache   map[string]*list.Element
	lruList *list.List
	maxSize int
}

type memoryEntry struct {
	key   string
	value *CacheEntry
}

// NewMemoryStorage creates a new in-memory cache with the specified maximum size.
// If maxSize is 0 or negative, a default size of 1000 entries is used.
func NewMemoryStorage(maxSize int) Cache {
	if maxSize <= 0 {
		maxSize = 1000
	}

	return &memoryStorage{
		cache:   make(map[string]*list.Element),
		lruList: list.New(),
		maxSize: maxSize,
	}
}

// Get retrieves an entry from the cache.
func (m *memoryStorage) Get(key string) (*CacheEntry, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	elem, ok := m.cache[key]
	if !ok {
		return nil, false
	}

	// Move to front (most recently used)
	m.lruList.MoveToFront(elem)

	entry := elem.Value.(*memoryEntry)
	return entry.value, true
}

// Set stores an entry in the cache.
func (m *memoryStorage) Set(key string, entry *CacheEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if key already exists
	if elem, ok := m.cache[key]; ok {
		// Update existing entry and move to front
		m.lruList.MoveToFront(elem)
		elem.Value.(*memoryEntry).value = entry
		return
	}

	// Add new entry
	elem := m.lruList.PushFront(&memoryEntry{
		key:   key,
		value: entry,
	})
	m.cache[key] = elem

	// Evict oldest if over capacity
	if m.lruList.Len() > m.maxSize {
		oldest := m.lruList.Back()
		if oldest != nil {
			m.lruList.Remove(oldest)
			oldEntry := oldest.Value.(*memoryEntry)
			delete(m.cache, oldEntry.key)
		}
	}
}

// Delete removes an entry from the cache.
func (m *memoryStorage) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if elem, ok := m.cache[key]; ok {
		m.lruList.Remove(elem)
		delete(m.cache, key)
	}
}

// Clear removes all entries from the cache.
func (m *memoryStorage) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache = make(map[string]*list.Element)
	m.lruList = list.New()
}
