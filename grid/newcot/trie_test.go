package newcot

import (
	"testing"

	"github.com/go-math/support/assert"
)

func TestTrieTap(t *testing.T) {
	trie := newTrie(3, 1)

	assert.Equal(trie.tap([]uint8{0, 0, 0}, []uint32{0, 1, 2}), false, t)
	assert.Equal(trie.tap([]uint8{0, 0, 0}, []uint32{0, 1, 2}), true, t)
	assert.Equal(trie.tap([]uint8{0, 0, 1}, []uint32{0, 1, 2}), false, t)
	assert.Equal(trie.tap([]uint8{0, 1, 0}, []uint32{2, 1, 0}), false, t)
	assert.Equal(trie.tap([]uint8{1, 0, 0}, []uint32{0, 2, 1}), false, t)
	assert.Equal(trie.tap([]uint8{1, 0, 0}, []uint32{0, 2, 1}), true, t)

	n := trie.root
	assert.Equal(len(n.children), 2, t)

	// 0
	n = n.children[0]
	assert.Equal(n.value, uint32(0), t)
	assert.Equal(len(n.children), 2, t)

	// 0, 0
	n = n.children[0]
	assert.Equal(n.value, uint32(0), t)
	assert.Equal(len(n.children), 2, t)

	// 0, 0, 0
	n = n.children[0]
	assert.Equal(n.value, uint32(0), t)
	assert.Equal(len(n.children), 1, t)

	// 0, 0, 0, 0
	n = n.children[0]
	assert.Equal(n.value, uint32(0), t)
	assert.Equal(len(n.children), 1, t)

	// 0, 0, 0, 0, 1
	n = n.children[0]
	assert.Equal(n.value, uint32(1), t)
	assert.Equal(len(n.children), 1, t)

	// 0, 0, 0, 0, 1, 2
	n = n.children[0]
	assert.Equal(n.value, uint32(2), t)
	assert.Equal(len(n.children), 0, t)
}


func TestTrieTap1D(t *testing.T) {
	t.Skip()
	trie := newTrie(1, 4)

	levels := []uint8{1, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4}
	orders := []uint32{0, 2, 1, 3, 1, 3, 5, 7, 1, 3, 5, 7, 9, 11, 13, 15}

	for i := range levels {
		assert.Equal(trie.tap(levels[i:i+1], orders[i:i+1]), false, t)
		assert.Equal(trie.tap(levels[i:i+1], orders[i:i+1]), true, t)
	}
}
