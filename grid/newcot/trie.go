package newcot

type trie struct {
	depth  uint16
	spread uint32
	root   *node
}

type node struct {
	value    uint32
	children []*node
}

// tap looks for the given sequence of numbers and returns true if found or
// returns false if not found, in which case the sequence is appended to the
// internal structure.
func (t *trie) tap(trace []uint32) bool {
	var i uint16
	var n, c *node

	n = t.root

outter:
	for i = 0; i < t.depth; i++ {
		for _, c = range n.children {
			if c.value == trace[i] {
				n = c
				continue outter
			}
		}

		// Not found.
		for ; i < t.depth; i++ {
			c = &node{
				value: trace[i],
				children: make([]*node, 0, t.spread),
			}
			n.children = append(n.children, c)
			n = c
		}

		return false
	}

	// Found.
	return true
}

func newTrie(depth uint16, spread uint32) *trie {
	return &trie{
		depth: depth,
		spread: spread,
		root: &node{
			children: make([]*node, 0, spread),
		},
	}
}
