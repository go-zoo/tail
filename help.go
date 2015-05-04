package tail

import (
	"io/ioutil"
	"log"

	"github.com/garyburd/redigo/redis"
)

func ReadTemplateFile(src string) ([]byte, error) {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Fetcher struct {
	Source redis.Conn
}

func NewFetcher(net string, addr string) *Fetcher {
	conn, err := redis.Dial(net, addr)
	if err != nil {
		return nil
	}
	return &Fetcher{conn}
}

func (f *Fetcher) GetData(id string) []byte {
	n, err := f.Source.Do("GET", id)
	if err != nil {
		log.Println(err)
		return nil
	}
	return n.([]byte)
}

func (f *Fetcher) SetData(id string, data interface{}) error {
	_, err := f.Source.Do("SET", id, data)
	if err != err {
		return err
	}
	return nil
}
