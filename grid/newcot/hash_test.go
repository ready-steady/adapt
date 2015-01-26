package newcot

import (
	"math/rand"
	"testing"
	"unsafe"

	"github.com/ready-steady/support/assert"
)

func TestHashTap(t *testing.T) {
	const (
		capacity = 10
	)

	hash := newHash(2, capacity)

	assert.Equal(hash.tap([]uint64{4, 2}), false, t)
	assert.Equal(hash.tap([]uint64{4, 2}), true, t)

	keys := make([]string, 0, capacity)
	for k, _ := range hash.mapping {
		keys = append(keys, k)
	}

	assert.Equal(len(keys), 1, t)

	if isLittleEndian() {
		assert.Equal(keys[0],
			"\x04\x00\x00\x00\x00\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00", t)
	} else {
		assert.Equal(keys[0],
			"\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00\x02", t)
	}
}

func isLittleEndian() bool {
	var x uint32 = 0x01020304
	return *(*byte)(unsafe.Pointer(&x)) == 0x04
}

func BenchmarkHashTap(b *testing.B) {
	const (
		dimensionCount = 20
		parentCount    = 1000

		depth    = dimensionCount
		capacity = 2 * parentCount * dimensionCount
	)

	generator := rand.New(rand.NewSource(0))
	data := make([]uint64, 2*parentCount*dimensionCount)
	for _, i := range data {
		data[i] = uint64(generator.Int63())
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hash := newHash(depth, capacity)
		for j := 0; j < parentCount; j++ {
			hash.tap(data[(2*j+0)*depth : (2*j+1)*depth])
			hash.tap(data[(2*j+1)*depth : (2*j+2)*depth])
		}
	}
}
