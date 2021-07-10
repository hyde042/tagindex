package tagindex_test

import (
	"testing"

	"github.com/hyde042/tagindex"
)

func TestCompare(t *testing.T) {
	assert(t, "1", tagindex.SetContains([]uint32{}, []uint32{}), true)
	assert(t, "2", tagindex.SetContains([]uint32{1, 2, 3, 4, 5}, []uint32{}), true)
	assert(t, "3", tagindex.SetContains([]uint32{1, 2, 3, 4, 5}, []uint32{1}), true)
	assert(t, "4", tagindex.SetContains([]uint32{1, 2, 3, 4, 5}, []uint32{2, 4}), true)
	assert(t, "5", tagindex.SetContains([]uint32{1, 2, 3, 4, 5}, []uint32{0, 4}), false)
	assert(t, "6", tagindex.SetContains([]uint32{1, 2, 3, 4, 5}, []uint32{2, 6}), false)
	assert(t, "7", tagindex.SetContains([]uint32{1, 2, 3, 4, 5}, []uint32{0}), false)
	assert(t, "8", tagindex.SetContains([]uint32{1, 2, 3, 4, 5}, []uint32{7}), false)
}
