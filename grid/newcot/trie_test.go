package newcot

import (
	"math/rand"
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestTrieTap(t *testing.T) {
	trie := newTrie(6, 3)

	assert.Equal(trie.tap([]uint64{0, 0, 0, 0, 1, 2}), false, t)
	assert.Equal(trie.tap([]uint64{0, 0, 0, 0, 1, 2}), true, t)
	assert.Equal(trie.tap([]uint64{0, 0, 1, 0, 1, 2}), false, t)
	assert.Equal(trie.tap([]uint64{0, 1, 0, 2, 1, 0}), false, t)
	assert.Equal(trie.tap([]uint64{1, 0, 0, 0, 2, 1}), false, t)
	assert.Equal(trie.tap([]uint64{1, 0, 0, 0, 2, 1}), true, t)

	n := trie.root
	assert.Equal(len(n.children), 2, t)

	// 0
	n = n.children[0]
	assert.Equal(n.value, uint64(0), t)
	assert.Equal(len(n.children), 2, t)

	// 0, 0
	n = n.children[0]
	assert.Equal(n.value, uint64(0), t)
	assert.Equal(len(n.children), 2, t)

	// 0, 0, 0
	n = n.children[0]
	assert.Equal(n.value, uint64(0), t)
	assert.Equal(len(n.children), 1, t)

	// 0, 0, 0, 0
	n = n.children[0]
	assert.Equal(n.value, uint64(0), t)
	assert.Equal(len(n.children), 1, t)

	// 0, 0, 0, 0, 1
	n = n.children[0]
	assert.Equal(n.value, uint64(1), t)
	assert.Equal(len(n.children), 1, t)

	// 0, 0, 0, 0, 1, 2
	n = n.children[0]
	assert.Equal(n.value, uint64(2), t)
	assert.Equal(len(n.children), 0, t)
}

func BenchmarkTrieTap(b *testing.B) {
	const (
		dimensionCount = 20
		parentCount    = 1000

		depth  = dimensionCount
		spread = 2 * parentCount * dimensionCount
	)

	generator := rand.New(rand.NewSource(0))
	data := make([]uint64, parentCount*depth)
	for _, i := range data {
		data[i] = uint64(generator.Int63())
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trie := newTrie(depth, spread)
		for j := 0; j < parentCount; j++ {
			trie.tap(data[j*depth : (j+1)*depth])
		}
	}
}
