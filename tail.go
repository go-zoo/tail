package tail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

type Asset struct {
	ID     string
	Source string
	TTL    time.Duration
	Data   []byte
	Cache  Cache
}

func New(id string, src string, cache Cache) *Asset {
	if cache != nil {
		asset := &Asset{ID: id, Source: src, Cache: cache}
		asset.Refresh()
		asset.WatchFile()
		return asset
	}
	log.Println("Cache arg cannot be nil")
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

func (a *Asset) WatchFile() {
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Println(err)
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						a.Refresh()
					}
				case err := <-watcher.Errors:
					log.Println(err)
				}

			}
		}()
		err = watcher.Add(a.Source)
		if err != nil {
			log.Println(err)

		}
		<-done
	}()
}

func (a *Asset) Refresh() {
	data, err := ReadAssetFile(a.Source)
	if err != nil {
		log.Println(err)
		return
	}
	err = a.Cache.Set(a.ID, data)
	if err != nil {
		log.Println(err)
	}
}

func (a *Asset) Get() []byte {
	return a.Cache.Get(a.ID)
}

func (a *Asset) Set(id string, data []byte) error {
	err := a.Cache.Set(a.ID, data)
	if err != nil {
		return err
	}
	return nil
}

func (a *Asset) Build() {
	var cmp []byte
	buffer := bytes.NewBuffer(cmp)

	tmp := template.New(a.ID)
	tmp.Parse(string(a.Cache.Get(a.ID)))
	tmp.Execute(buffer, a.Data)

	a.Set(a.ID, buffer.Bytes())

	fmt.Println(buffer.String())
}
