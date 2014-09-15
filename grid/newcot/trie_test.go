package newcot

import (
	"testing"

	"github.com/go-math/support/assert"
)

func TestTrieTap(t *testing.T) {
	const (
		depth  = 5
		spread = 5
	)

	trie := newTrie(depth, spread)

	assert.Equal(trie.tap([]uint32{0, 1, 2, 3, 4, 5}), false, t)
	assert.Equal(trie.tap([]uint32{0, 1, 2, 3, 4, 5}), true, t)
	assert.Equal(trie.tap([]uint32{0, 1, 3, 3, 3, 3}), false, t)
	assert.Equal(trie.tap([]uint32{0, 1, 3, 3, 4, 4}), false, t)
	assert.Equal(trie.tap([]uint32{0, 1, 3, 3, 4, 4}), true, t)

	n := trie.root
	assert.Equal(len(n.children), 1, t)

	// 0
	n = n.children[0]
	assert.Equal(len(n.children), 1, t)

	// 0, 1
	n = n.children[0]
	assert.Equal(len(n.children), 2, t)

	// 0, 1, 2
	m := n.children[0]
	assert.Equal(len(m.children), 1, t)

	// 0, 1, 3
	m = n.children[1]
	assert.Equal(len(m.children), 1, t)

	// 0, 1, 3, 3
	n = m.children[0]
	assert.Equal(len(n.children), 2, t)
}
