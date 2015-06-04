package tail

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	_ "github.com/boltdb/bolt"
	_ "github.com/garyburd/redigo/redis"
)

type Asset struct {
	ID      string
	Source  string
	Data    []byte
	TTL     time.Duration
	Clients map[string]*Client
	Cache   Cache
}

func New(id string, src string, ttl time.Duration, cache Cache) (*Asset, error) {
	if cache != nil {
		asset := &Asset{
			ID: id, Source: src,
			TTL:     ttl,
			Clients: make(map[string]*Client),
			Cache:   cache,
		}
		asset.loadAsset()
		asset.WatchFile()
		asset.cleanClient()
		return asset, nil
	}
	return nil, errors.New("Cache cannot be nil")
}

func (a *Asset) cleanClient() {
	go func() {
		for {
			time.AfterFunc(a.TTL, func() {
				for _, c := range a.Clients {
					if !c.Valid {
						fmt.Printf("Client: %s has expired.\n", c.ID)
						err := a.Del(c.ID)
						fmt.Println(c.ID, "Deleted")
						if err != nil {
							fmt.Println(err)
						}
					}

				}
			})
			time.Sleep(a.TTL)
		}
	}()
}

func (a *Asset) loadAsset() error {
	var err error
	a.Data, err = ReadAssetFile(a.Source)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) Get(id string) []byte {
	cid := a.buildId(id)
	if a.Clients[cid] == nil {
		return nil
	}
	a.Clients[cid].renewClient()
	return a.Cache.Get(cid)
}

func (a *Asset) Set(id string, data []byte) error {
	cid := a.buildId(id)
	if a.Clients[cid] == nil {
		clt := NewClient(cid, a.TTL)
		a.Clients[cid] = clt
	}
	err := a.Cache.Set(cid, data)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) Update(id string, data []byte) error {
	err := a.Cache.Update(id, data)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) Del(id string) error {
	delete(a.Clients, id)
	err := a.Cache.Del(id)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) create(id string, data interface{}) error {
	var buff bytes.Buffer
	tmp := template.New(a.ID)
	tmp.Parse(string(a.Data))
	err := tmp.Execute(&buff, data)
	if err != nil {
		return err
	}
	err = a.Set(id, buff.Bytes())
	if err != nil {
		return err
	}

	cid := a.buildId(id)
	clt := NewClient(cid, a.TTL)
	a.Clients[cid] = clt
	return nil
}

func (a *Asset) Create(id string, data interface{}) error {
	err := a.create(id, data)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) GetOrNew(id string, data interface{}) ([]byte, error) {
	cid := a.buildId(id)
	if a.Clients[cid] == nil {
		err := a.create(id, data)
		if err != nil {
			return nil, err
		}
	}
	a.Clients[cid].renewClient()
	return a.Cache.Get(cid), nil
}

func (a *Asset) buildId(id string) string {
	return fmt.Sprintf("%s:%s", a.ID, id)
}
