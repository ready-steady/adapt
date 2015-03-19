package adhier

import (
	"math/rand"
	"sort"
	"testing"
	"unsafe"

	"github.com/ready-steady/support/assert"
)

func TestHashUnique(t *testing.T) {
	hash := newHash(2)

	assert.Equal(hash.unique([]uint64{4, 2}), []uint64{4, 2}, t)
	assert.Equal(hash.unique([]uint64{6, 9}), []uint64{6, 9}, t)
	assert.Equal(hash.unique([]uint64{4, 2}), []uint64{}, t)

	keys := make([]string, 0)
	for k, _ := range hash.mapping {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))

	assert.Equal(len(keys), 2, t)

	if isLittleEndian() {
		assert.Equal(keys[0],
			"\x04\x00\x00\x00\x00\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00", t)
		assert.Equal(keys[1],
			"\x06\x00\x00\x00\x00\x00\x00\x00\x09\x00\x00\x00\x00\x00\x00\x00", t)
	} else {
		assert.Equal(keys[0],
			"\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00\x02", t)
		assert.Equal(keys[1],
			"\x00\x00\x00\x00\x00\x00\x00\x06\x00\x00\x00\x00\x00\x00\x00\x09", t)
	}
}

func TestHashUniqueOverlap(t *testing.T) {
	hash := newHash(2)

	key := []uint64{4, 2}
	assert.Equal(hash.unique(key), key, t)

	key[0], key[1] = 6, 9
	key = []uint64{4, 2}
	assert.Equal(hash.unique(key), []uint64{}, t)
}

func isLittleEndian() bool {
	var x uint32 = 0x01020304
	return *(*byte)(unsafe.Pointer(&x)) == 0x04
}

func BenchmarkHashUnique(b *testing.B) {
	const (
		ni = 20
		nn = 1000
	)

	hash := newHash(ni)

	generator := rand.New(rand.NewSource(0))
	indices := make([]uint64, 2*nn*ni)
	for _, i := range indices {
		indices[i] = uint64(generator.Int63())
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hash.unique(indices)
	}
}
