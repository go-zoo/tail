package rediscache

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// Redis Assets caching

type RedisCache struct {
	source redis.Conn
}

func New(net string, addr string) *RedisCache {
	conn, err := redis.Dial(net, addr)
	if err != nil {
		return nil
	}
	return &RedisCache{conn}
}

func (f *RedisCache) Get(id string) []byte {
	n, err := f.source.Do("GET", id)
	if err != nil {
		log.Println(err)
		return nil
	}
	if n != nil {
		return n.([]byte)
	}
	return nil
}

func (f *RedisCache) Set(id string, data []byte) error {
	_, err := f.source.Do("SET", id, data)
	if err != err {
		return err
	}
	return nil
}
