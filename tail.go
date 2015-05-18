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
	Data    bytes.Buffer
	TTL     time.Duration
	Clients map[string]*Client
	Cache   Cache
}

type Client struct {
	ID     string
	Valid  bool
	TTL    time.Duration
	Expire *time.Timer
}

func New(id string, src string, ttl int64, cache Cache) *Asset {
	if cache != nil {
		asset := &Asset{
			ID: id, Source: src,
			TTL:     time.Duration(ttl),
			Clients: make(map[string]*Client),
			Cache:   cache,
		}

		asset.Refresh()
		asset.WatchFile()
		go asset.cleanClient()
		return asset
	}
	fmt.Println("Cache arg cannot be nil")
	return nil
}

func NewClient(id string, ttl time.Duration) *Client {
	c := &Client{ID: id, Valid: true, TTL: ttl, Expire: time.NewTimer(ttl)}
	go c.watch()
	return c
}

func (a *Asset) cleanClient() {
	for {
		time.AfterFunc(a.TTL, func() {
			for _, c := range a.Clients {
				if !c.Valid {
					fmt.Printf("Client: %s have expire.\n", c.ID)
					err := a.Del(c.ID)
					if err != nil {
						fmt.Println(err)
					}
				}

			}
		})
		time.Sleep(a.TTL)
	}

}

func (a *Asset) Refresh() {
	data, err := ReadAssetFile(a.Source)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = a.Cache.Set(a.ID, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (a *Asset) del(id string) error {
	err := a.Cache.Del(id)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) Del(id string) error {
	delete(a.Clients, id)
	return a.del(id)
}

func (a *Asset) get() []byte {
	return a.Cache.Get(a.ID)
}

func (a *Asset) Get(id string) ([]byte, error) {
	cid := a.buildId(id)

	a.Clients[cid].renewClient()
	d := a.Cache.Get(cid)
	if d == nil {
		return nil, errors.New(fmt.Sprintf("%s not found !", id))
	}
	return d, nil
}

func (c *Client) renewClient() {
	c.Expire.Reset(c.TTL)
	//c.Expire = time.NewTimer(c.TTL)
	c.Valid = true
}

func (a *Asset) GetOrBuild(id string, data interface{}) ([]byte, error) {
	cid := a.buildId(id)
	if a.Clients[cid] == nil {
		err := a.build(id, data)
		if err != nil {
			return nil, err
		}
	}
	a.Clients[cid].renewClient()
	return a.Cache.Get(cid), nil
}

func (a *Asset) Set(id string, data []byte) error {
	err := a.Cache.Set(id, data)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) set(id string, data []byte) error {
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

func (c *Client) watch() {
	<-c.Expire.C
	c.Valid = false
}

func (a *Asset) build(id string, data interface{}) error {
	tmp := template.New(a.ID)
	tmp.Parse(string(a.Cache.Get(a.ID)))
	tmp.Execute(&a.Data, data)
	err := a.set(id, a.Data.Bytes())
	if err != nil {
		return err
	}
	a.Data.Reset()
	cid := a.buildId(id)
	clt := NewClient(cid, a.TTL)
	a.Clients[cid] = clt
	return nil
}

func (a *Asset) Build(id string, data interface{}) error {
	err := a.build(id, data)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) buildId(id string) string {
	return fmt.Sprintf("%s:%s", a.ID, id)
}
