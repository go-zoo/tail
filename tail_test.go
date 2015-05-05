package tail

import (
	"testing"

	"github.com/go-zoo/tail/memcache"
)

func TestTemplateCreation(t *testing.T) {
	tpl := New("test", "", nil)
	if tpl == nil {
		t.Fail()
	}
}

func TestTemplateCaching(t *testing.T) {
	cache := memcache.New()
	tpl := New("test", "tail_test.go", cache)
	if tpl.Get() == nil {
		t.Fail()
	}
}
