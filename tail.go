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
	ID     string
	Source string
	Data   bytes.Buffer
	Cache  Cache
}

func New(id string, src string, cache Cache) *Asset {
	if cache != nil {
		asset := &Asset{ID: id, Source: src, Cache: cache}
		asset.Refresh()
		asset.WatchFile()
		return asset
	}
	fmt.Println("Cache arg cannot be nil")
	return nil
}

func (a *Asset) Watch(ttl time.Duration) {
	fmt.Printf("[+] Refreshing %s Asset each %s\n", a.Source, ttl.String())
	go func() {
		for {
			time.AfterFunc(ttl, a.Refresh)
			time.Sleep(ttl)
		}
	}()
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

func (a *Asset) get() []byte {
	return a.Cache.Get(a.ID)
}

func (a *Asset) Get(id string) ([]byte, error) {
	d := a.Cache.Get(a.buildId(id))
	if d == nil {
		return nil, errors.New(fmt.Sprintf("%s not found !", id))
	}
	return d, nil
}

func (a *Asset) GetOrBuild(id string, data interface{}) ([]byte, error) {
	d := a.Cache.Get(a.buildId(id))
	if d == nil {
		err := a.build(id, data)
		if err != nil {
			return nil, err
		}
		return a.Cache.Get(a.buildId(id)), nil
	}
	return d, nil
}

func (a *Asset) Set(id string, data []byte) error {
	err := a.Cache.Set(id, data)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) set(id string, data []byte) error {
	err := a.Cache.Set(a.buildId(id), data)
	if err != nil {
		return err
	}
	return nil
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
