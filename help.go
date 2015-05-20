package tail

import (
	"fmt"
	"io/ioutil"

	"github.com/go-fsnotify/fsnotify"
)

func ReadAssetFile(src string) ([]byte, error) {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *Asset) WatchFile() {
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println(err)
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						a.loadAsset()
					}
				case err := <-watcher.Errors:
					fmt.Println(err)
				}

			}
		}()
		err = watcher.Add(a.Source)
		if err != nil {
			fmt.Println(err)

		}
		<-done
	}()
}
