package tagindex

import (
	"encoding/binary"

	"github.com/cespare/xxhash"
)

type bloom [4]uint64

func makeBloom(data []uint32, k int) bloom {
	var (
		h = xxhash.New() // fnv.New64a()
		b [4]byte
		f bloom
	)
	for i, n := range data {
		if i > 0 {
			h.Reset()
		}
		for j := 0; j < k; j++ {
			binary.LittleEndian.PutUint32(b[:], n)
			h.Write(b[:])
			var (
				a = h.Sum64() % 256
				b = a / 64
			)
			f[b] = f[b] | 1<<(a-b*64)
		}
	}
	return f
}

func (t bloom) contains(o bloom) bool {
	return t[0]&o[0] == o[0] && t[1]&o[1] == o[1] && t[2]&o[2] == o[2] && t[3]&o[3] == o[3]
}
