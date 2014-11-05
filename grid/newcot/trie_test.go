package newcot

import (
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
