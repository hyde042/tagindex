package tagindex_test

import (
	"testing"

	"github.com/hyde042/tagindex"
)

func TestFilter(t *testing.T) {
	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 1))
	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 2))
	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 3))
	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 10))
	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 100))

	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 9).Contains(tagindex.MakeBloom([]uint32{}, 9)))
	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 9).Contains(tagindex.MakeBloom([]uint32{0}, 9)))
	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 9).Contains(tagindex.MakeBloom([]uint32{0, 64}, 9)))
	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 9).Contains(tagindex.MakeBloom([]uint32{100, 150}, 9)))
	t.Log(tagindex.MakeBloom([]uint32{0, 63, 64, 255, 257}, 9).Contains(tagindex.MakeBloom([]uint32{1}, 9)))
}
