package tail

import (
	"testing"

	"github.com/go-zoo/tail/memcache"
)

func TestTemplateCreation(t *testing.T) {
	cache := memcache.New()
	tpl, _ := New("test", "help.go", 5, cache)
	if tpl == nil {
		t.Fail()
	}
}

func TestTemplateCaching(t *testing.T) {
	cache := memcache.New()
	tpl, _ := New("test", "tail_test.go", 5, cache)
	if tpl.get() == nil {
		t.Fail()
	}
}

func TestMultipleTemplate(t *testing.T) {
	cache := memcache.New()
	tpl1, _ := New("1", "tail_test.go", 5, cache)
	tpl2, _ := New("2", "tail.go", 5, cache)
	if tpl1.get() == nil || tpl2.get() == nil {
		t.Fail()
	}
}

func TestMultipleCache(t *testing.T) {
	c1 := memcache.New()
	c2 := memcache.New()
	tpl1, _ := New("1", "tail_test.go", 5, c1)
	tpl2, _ := New("2", "tail.go", 5, c2)
	if tpl1.get() == nil || tpl2.get() == nil {
		t.Fail()
	}
}
