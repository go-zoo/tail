package main

import (
	"net/http"

	"github.com/go-zoo/tail"
)

var (
	fetcher   = tail.NewFetcher("tcp", "104.236.16.169:6379")
	IndexTmpl = tail.New("index", "index.html", fetcher)
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
