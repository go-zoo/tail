tail [![GoDoc](https://godoc.org/github.com/go-zoo/tail?status.png)](http://godoc.org/github.com/go-zoo/tail) [![Build Status](https://travis-ci.org/go-zoo/tail.svg)](https://travis-ci.org/go-zoo/tail)
=======
### Work in progress !

## What is tail ?

Tail is a easy to use caching system for your web project. 
It include multiple caching platform {redis, bolt, in memory} and more will come later.
Once the ressource is in cache tail will update it if any change append.
You can also force update at each X secondes with `indexTemplate.Watch("5s")`.

![alt tag](http://blackheartmagazine.com/blog/wp-content/uploads/2014/01/tail.jpg)

## Example

``` go

package main

import(
  "net/http"

  // Main package for using tail
  "github.com/go-zoo/tail"
  // import depend of the platform you choose
  "github.com/go-zoo/tail/rediscache"
)

var (
  // Create the cache platform with your parameters
  cache = rediscache.New("tcp", "localhost:6379")

  // Create a new template that takes (id, path, cachePlatform)
  indexTemplate = tail.New("index", "template/index.html", cache)
)

func main () {
  http.HandleFunc("/", indexHandler)

  http.ListenAndServe(":8080", mux)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
  rw.Write(indexTemplate.Get())
}

```
## TODO

- More Functionnalities
- DOC
- More Testing
- Debugging
- Optimisation

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Write Tests!
4. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request

## License
MIT

## Links

Lightning fast http Multiplexer : [Bone](https://github.com/go-zoo/bone)
