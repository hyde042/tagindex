package tagindex_test

import (
	"testing"

	"github.com/hyde042/tagindex"
)

func TestBloom(t *testing.T) {
	t.Log(makeBloom(1, 0, 63, 64, 255, 257))
	t.Log(makeBloom(2, 0, 63, 64, 255, 257))
	t.Log(makeBloom(3, 0, 63, 64, 255, 257))
	t.Log(makeBloom(10, 0, 63, 64, 255, 257))
	t.Log(makeBloom(100, 0, 63, 64, 255, 257))

	testBloom := makeBloom(9, 0, 63, 64, 255, 257)
	assert(t, "1", testBloom.Contains(makeBloom(9)), true)
	assert(t, "2", testBloom.Contains(makeBloom(9, 0)), true)
	assert(t, "3", testBloom.Contains(makeBloom(9, 0, 64)), true)
	assert(t, "4", testBloom.Contains(makeBloom(9, 100, 150)), false)
	assert(t, "5", testBloom.Contains(makeBloom(9, 1)), false)
}

func makeBloom(k int, values ...uint32) tagindex.Bloom {
	return tagindex.MakeBloom(values, k)
}
