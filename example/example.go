package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-zoo/tail"
	"github.com/go-zoo/tail/rediscache"
)

type Data struct {
	Name string `json:"name"`
}

var (
	cache = rediscache.New("tcp", "104.236.16.169:6379")
	//memcache.New()
	//rediscache.New("tcp", "104.236.16.169:6379")
	//boltcache.New("fetch.db", 0600, nil)

	IndexTmpl = tail.New("index", "index.html", cache)
	Img       = tail.New("img", "logo.png", cache)
)

func init() {
	n := Data{"Some APP"}
	raw, _ := json.Marshal(n)
	IndexTmpl.Data = raw
	IndexTmpl.Build()
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/img", imgHandler)

	http.ListenAndServe(":8080", nil)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Write(IndexTmpl.Get())
}

func imgHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Write(Img.Get())
}
