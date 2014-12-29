package newcot

type trie struct {
	depth  uint32
	spread uint32
	root   *node
}

type node struct {
	value    uint64
	children []*node
}

func newTrie(depth uint32, spread uint32) *trie {
	return &trie{
		depth:  depth,
		spread: spread,
		root: &node{
			children: make([]*node, 0, spread),
		},
	}
}

// tap looks for the given sequence of numbers and returns true if found or
// returns false if not found, in which case the sequence is appended to the
// internal structure.
func (t *trie) tap(trace []uint64) bool {
	var c *node

outer:
	for n, i := t.root, uint32(0); i < t.depth; i++ {
		j := uint32(0)
		count := uint32(len(n.children))

		for ; j < count; j++ {
			c = n.children[j]
			if c.value == trace[i] {
				n = c
				continue outer
			} else if c.value > trace[i] {
				// The children are always kept sorted.
				break
			}
		}

		k := int32(t.depth - 1)

		// Create a leaf.
		c = &node{
			value: trace[k],
		}

		// Create the rest of the tail.
		for k--; k >= int32(i); k-- {
			children := make([]*node, 1, t.spread)
			children[0] = c
			c = &node{
				value:    trace[k],
				children: children,
			}
		}

		// Insert the node at the jth position.
		n.children = n.children[:count+1]
		copy(n.children[j+1:], n.children[j:count])
		n.children[j] = c

		return false
	}

	return true
}
