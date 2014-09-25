package newcot

// Closed represents an instance of the grid on [0, 1]^n.
type Closed struct {
	dc uint16
}

// NewClosed creates an instance of the grid on [0, 1]^n.
func NewClosed(dimensions uint16) *Closed {
	return &Closed{dimensions}
}

// Dimensions returns the number of dimensions of the grid.
func (c *Closed) Dimensions() uint16 {
	return c.dc
}

// ComputeNodes returns the nodes corresponding to the given index.
func (_ *Closed) ComputeNodes(index []uint64) []float64 {
	nodes := make([]float64, len(index))

	for i := range nodes {
		if uint32(index[i]) == 0 {
			nodes[i] = 0.5
		} else {
			nodes[i] = float64(index[i]>>32) / float64(uint32(2)<<(uint32(index[i])-1))
		}
	}

	return nodes
}

// ComputeChildren returns the index of the child nodes corresponding to the
// parent nodes given by their index.
func (c *Closed) ComputeChildren(parentIndex []uint64) []uint64 {
	dc := uint32(c.dc)
	pc := uint32(len(parentIndex)) / dc

	index := make([]uint64, 2*pc*dc*dc)

	// The algorithm needs to keep track and eliminate duplicate nodes. To this
	// end, a trie (https://en.wikipedia.org/wiki/Trie) is utilized.
	trie := newTrie(dc, 2*pc*dc)

	cc := uint32(0)

	push := func(p, d uint32, pair uint64) {
		copy(index[cc*dc:], parentIndex[p*dc:(p+1)*dc])
		index[cc*dc+d] = pair

		if !trie.tap(index[cc*dc:]) {
			cc++
		}
	}

	var i, j, level, order uint32

	for i = 0; i < pc; i++ {
		for j = 0; j < dc; j++ {
			level = uint32(parentIndex[i*dc+j])

			if level == 0 {
				push(i, j, 1|0<<32)
				push(i, j, 1|2<<32)
				continue
			}

			order = uint32(parentIndex[i*dc+j] >> 32)

			if level == 1 {
				push(i, j, 2|uint64(order+1)<<32)
			} else {
				push(i, j, uint64(level+1)|uint64(2*order-1)<<32)
				push(i, j, uint64(level+1)|uint64(2*order+1)<<32)
			}
		}
	}

	return index[0 : cc*dc]
}
