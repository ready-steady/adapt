package newcot

type trie struct {
	ic   uint32
	ml   uint32
	root *node
}

type node struct {
	value    uint32
	children []*node
}

func newTrie(ic uint16, ml uint8) *trie {
	return &trie{
		ic: uint32(ic),
		ml: uint32(ml),
		root: &node{
			children: make([]*node, 0, ml+1),
		},
	}
}

// tap looks for the given sequence of levels and orders and returns true if
// found or returns false if not found, in which case the sequence is appended
// to the internal structure.
func (t *trie) tap(levels []uint8, orders []uint32) bool {
	var i, j, value, count uint32
	var k int32
	var c, g *node

	n := t.root

overLevels:
	for i = 0; i < t.ic; i++ {
		value = uint32(levels[i])
		count = uint32(len(n.children))

		for j = 0; j < count; j++ {
			c = n.children[j]
			if c.value == value {
				n = c
				continue overLevels
			} else if c.value > value {
				// The children are always kept sorted.
				break
			}
		}

		// Create an order leaf.
		g = &node{
			value: orders[t.ic - 1],
		}

		// Create the rest of order nodes.
		for k = int32(t.ic) - 2; k >= 0; k-- {
			c = &node{
				value:    orders[k],
				children: make([]*node, 1, 1<<t.ml+1),
			}
			c.children[0] = g
			g = c
		}

		if i == t.ic - 1 {
			// Create the missing level node.
			c = &node{
				value:    uint32(levels[i]),
				children: make([]*node, 1, 1<<t.ml+1),
			}
			c.children[0] = g
		} else {
			// Create a level leaf.
			c = &node{
				value:    uint32(levels[t.ic - 1]),
				children: make([]*node, 1, 1<<t.ml+1),
			}
			c.children[0] = g
			g = c

			// Create the rest of level nodes down to i exclusive.
			for k = int32(t.ic) - 2; k > int32(i); k-- {
				c = &node{
					value:    uint32(levels[k]),
					children: make([]*node, 1, t.ml+1),
				}
				c.children[0] = g
				g = c
			}

			// Create the missing level node.
			c = &node{
				value:    uint32(levels[i]),
				children: make([]*node, 1, t.ml+1),
			}
			c.children[0] = g
		}

		// Insert the node at the jth position.
		n.children = n.children[:count+1]
		copy(n.children[j+1:], n.children[j:count])
		n.children[j] = c

		return false
	}

overOrders:
	for i = 0; i < t.ic; i++ {
		value = orders[i]
		count = uint32(len(n.children))

		for j = 0; j < count; j++ {
			c = n.children[j]
			if c.value == value {
				n = c
				continue overOrders
			} else if c.value > value {
				// The children are always kept sorted.
				break
			}
		}

		if i == t.ic - 1 {
			// Create the missing order node.
			c = &node{
				value: orders[i],
			}
		} else {
			// Create an order leaf.
			g = &node{
				value: orders[t.ic - 1],
			}

			// Create the rest of order nodes up to i exclusive.
			for k = int32(t.ic) - 2; k > int32(i); k-- {
				c = &node{
					value:    orders[k],
					children: make([]*node, 1, 1<<t.ml+1),
				}
				c.children[0] = g
				g = c
			}

			// Create the missing order node.
			c = &node{
				value:    orders[i],
				children: make([]*node, 1, 1<<t.ml+1),
			}
			c.children[0] = g
		}

		// Insert the node at the jth position.
		n.children = n.children[:count+1]
		copy(n.children[j+1:], n.children[j:count])
		n.children[j] = c

		return false
	}

	return true
}
