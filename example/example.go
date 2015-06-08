package main

import (
	"net/http"
	"time"

	"github.com/go-zoo/tail"
	"github.com/go-zoo/tail/memcache"
)

type Data struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	cache = memcache.New()
	//memcache.New()
	//rediscache.New("tcp", "redis-cache-1.squiidz.cont.tutum.io:49153")
	//boltcache.New("fetch.db", 0600, nil)

	IndexTmpl, _ = tail.New("index", "index.html", time.Second*5, cache)
	Img, _       = tail.New("img", "logo.png", time.Second*10, cache)
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/img", imgHandler)

	http.ListenAndServe(":9000", nil)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	data, _ := IndexTmpl.GetOrNew(req.RemoteAddr, &Data{req.RequestURI, 24})
	rw.Write(data)
}

func imgHandler(rw http.ResponseWriter, req *http.Request) {
	data := Img.Get("")
	rw.Write(data)
}
