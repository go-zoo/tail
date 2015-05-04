package tail

import (
	"fmt"
	"log"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

type Template struct {
	ID      string
	Source  string
	TTL     time.Duration
	Content Fetch
}

func New(id string, src string, ftch Fetch) *Template {
	tmpl := &Template{ID: id, Source: src, Content: ftch}
	tmpl.Refresh()
	return tmpl
}

func (t *Template) Watch(ttl time.Duration) {
	fmt.Printf("[+] Refreshing %s template each %s\n", t.Source, ttl.String())
	go func() {
		for {
			time.AfterFunc(ttl, t.Refresh)
			time.Sleep(ttl)
		}
	}()
}

func (t *Template) WatchFile() {
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
						t.Refresh()
					}
				case err := <-watcher.Errors:
					log.Println(err)
				}

			}
		}()
		err = watcher.Add(t.Source)
		if err != nil {
			log.Println(err)

		}
		<-done
	}()
}

func (t *Template) Refresh() {
	data, err := ReadTemplateFile(t.Source)
	if err != nil {
		log.Println(err)
	}
	err = t.Content.SetData(t.ID, data)
	if err != nil {
		log.Println(err)
	}
}

func (t *Template) Get() []byte {
	return t.Content.GetData(t.ID)
}

func (t *Template) Set(id string, data interface{}) error {
	err := t.Content.SetData(t.ID, data)
	if err != nil {
		return err
	}
	return nil
}
