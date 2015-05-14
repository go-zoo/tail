package main

import (
	"net/http"
	"time"

	"github.com/go-zoo/tail"
	"github.com/go-zoo/tail/rediscache"
)

type Data struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	cache, _ = rediscache.New("tcp", "104.236.16.169:6379")
	//memcache.New()
	//rediscache.New("tcp", "104.236.16.169:6379")
	//boltcache.New("fetch.db", 0600, nil)

	IndexTmpl = tail.New("index", "index.html", cache)
	Img       = tail.New("img", "logo.png", cache)
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/img", imgHandler)

	http.ListenAndServe(":8080", nil)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		d := Data{Name: time.Now().String()}
		IndexTmpl.Rebuild("911205", d)
	}
	data, _ := IndexTmpl.Get("911205")
	rw.Write(data)
}

func imgHandler(rw http.ResponseWriter, req *http.Request) {
	data, _ := Img.Get("")
	rw.Write(data)
}
