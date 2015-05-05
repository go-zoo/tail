package boltcache

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
)

// BoltDB template caching

type BoltCache struct {
	source *bolt.DB
}

func New(path string, mode os.FileMode, options *bolt.Options) *BoltCache {
	boltC, err := bolt.Open(path, mode, options)
	if err != nil {
		log.Println(err)
	}
	return &BoltCache{boltC}
}

func (b *BoltCache) Get(id string) []byte {
	var val []byte
	b.source.View(func(tx *bolt.Tx) error {
		val = tx.Bucket([]byte("default")).Get([]byte(id))
		return nil
	})
	return val
}

func (b *BoltCache) Set(id string, data []byte) error {
	b.source.Update(func(tx *bolt.Tx) error {
		buck, _ := tx.CreateBucketIfNotExists([]byte("default"))
		buck.Put([]byte(id), data)
		return nil
	})
	return nil
}
