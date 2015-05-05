package main

import (
	"net/http"

	"github.com/go-zoo/tail"
	"github.com/go-zoo/tail/memcache"
)

var (
	cache = memcache.New()
	//rediscache.New("tcp", "104.236.16.169:6379")
	//boltcache.New("fetch.db", 0600, nil)

	IndexTmpl = tail.New("index", "index.html", cache)
)

func init() {
	IndexTmpl.WatchFile()
}

func main() {
	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8080", nil)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Write(IndexTmpl.Get())
}
