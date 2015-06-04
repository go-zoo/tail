package memcache

// In Memory template caching

import (
	"errors"
	"fmt"
)

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
	if m.source != nil {
		m.source[id] = data
		return nil
	}
	return errors.New("Source is nil")
}

func (m *MemoryCache) Update(id string, data []byte) error {
	if m.source[id] != nil {
		m.source[id] = data
		return nil
	}
	return errors.New(fmt.Sprintf("%s doesn't exist\n", id))
}

func (m *MemoryCache) Del(id string) error {
	if m.source[id] == nil {
		return errors.New("Key doesn't exist")
	}
	delete(m.source, id)
	return nil
}
