package tagindex_test

import (
	"testing"

	"github.com/hyde042/tagindex"
)

func TestIndex(t *testing.T) {
	idx := tagindex.New()
	idx.Put(tagindex.Entry{ID: "1", Tags: []string{"a"}})
	idx.Put(tagindex.Entry{ID: "4", Tags: []string{"b", "d"}})
	idx.Put(tagindex.Entry{ID: "2", Tags: []string{"a", "b", "c"}})
	idx.Put(tagindex.Entry{ID: "3", Tags: []string{"b", "c"}})

	assert(t, "1", idx.Query([]string{"a"}, 10).TotalCount == 2, true)
	assert(t, "2", idx.Query([]string{"b"}, 10).TotalCount == 3, true)

}

func assert(t *testing.T, tag string, result, expected bool) {
	if result != expected {
		t.Fatalf("%s: result mismatch", tag)
	}
}
