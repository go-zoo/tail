package tail

import (
	"testing"

	"github.com/go-zoo/tail/memcache"
)

func TestTemplateCreation(t *testing.T) {
	cache, _ := memcache.New()
	_, err := New("test", "help.go", 5, cache)
	if err != nil {
		t.Fail()
	}
}

func TestTemplateCaching(t *testing.T) {
	cache, _ := memcache.New()
	_, err := New("test", "tail_test.go", 5, cache)
	if err != nil {
		t.Fail()
	}
}

func TestMultipleTemplate(t *testing.T) {
	cache, _ := memcache.New()
	_, err := New("1", "tail_test.go", 5, cache)
	_, err2 := New("2", "tail.go", 5, cache)
	if err != nil || err2 != nil {
		t.Fail()
	}
}

func TestMultipleCache(t *testing.T) {
	c1, _ := memcache.New()
	c2, _ := memcache.New()
	_, err := New("1", "tail_test.go", 5, c1)
	_, err2 := New("2", "tail.go", 5, c2)
	if err != nil || err2 != nil {
		t.Fail()
	}
}
