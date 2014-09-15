package newcot

type trie struct {
	ic   uint16
	ml   uint8
	root *node
}

type node struct {
	value    uint32
	children []*node
}

// tap looks for the given sequence of levels and orders and returns true if
// found or returns false if not found, in which case the sequence is appended
// to the internal structure.
func (t *trie) tap(levels []uint8, orders []uint32) bool {
	var i uint16
	var value uint32
	var n, c *node

	n = t.root

overLevels:
	for i = 0; i < t.ic; i++ {
		value = uint32(levels[i])

		for _, c = range n.children {
			if c.value == value {
				n = c
				continue overLevels
			}
		}

		for ; i < t.ic; i++ {
			c = &node{
				value:    uint32(levels[i]),
				children: make([]*node, 0, t.ml+1),
			}
			n.children = append(n.children, c)
			n = c
		}
	}

overOrders:
	for i = 0; i < t.ic; i++ {
		value = orders[i]

		for _, c = range n.children {
			if c.value == value {
				n = c
				continue overOrders
			}
		}

		for ; i < t.ic; i++ {
			c = &node{
				value:    orders[i],
				children: make([]*node, 0, 1<<t.ml+1),
			}
			n.children = append(n.children, c)
			n = c
		}

		return false
	}

	return true
}

func newTrie(ic uint16, ml uint8) *trie {
	return &trie{
		ic: ic,
		ml: ml,
		root: &node{
			children: make([]*node, 0, ml+1),
		},
	}
}
