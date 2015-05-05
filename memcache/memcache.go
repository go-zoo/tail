package memcache

// In Memory template caching

type MemoryCache struct {
	source map[string][]byte
}

func New() *MemoryCache {
	return &MemoryCache{source: make(map[string][]byte)}
}

func (m *MemoryCache) Get(id string) []byte {
	return m.source[id]
}

func (m *MemoryCache) Set(id string, data []byte) error {
	m.source[id] = data
	return nil
}
