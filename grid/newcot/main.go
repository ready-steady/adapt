// Package newcot provides means for working with the Newton–Cotes grid on a
// unit hypercube including boundaries.
package newcot

// Self represents a particular instantiation of the grid.
type Self struct {
	dc uint16
}

// New creates an instance of the grid for the given dimensionality.
func New(dimensions uint16) *Self {
	return &Self{dimensions}
}

// Dimensions returns the number of dimensions of the grid.
func (self *Self) Dimensions() uint16 {
	return self.dc
}

// ComputeNodes returns the nodes corresponding to the given index, which is
// a compact representation of a sequence of pairs (level, order).
func (_ *Self) ComputeNodes(index []uint64) []float64 {
	nodes := make([]float64, len(index))

	for i := range nodes {
		level := uint8(index[i] >> 32)
		if level == 0 {
			nodes[i] = 0.5
		} else {
			nodes[i] = float64(uint32(index[i])) / float64(uint32(2)<<(level-1))
		}
	}

	return nodes
}

// ComputeChildren returns the index of the child nodes corresponding to the
// parent nodes given by their index.
func (self *Self) ComputeChildren(parentIndex []uint64) []uint64 {
	dc := uint32(self.dc)
	pc := uint32(len(parentIndex)) / dc

	index := make([]uint64, 2*pc*dc*dc)

	// Create a trie for keeping track of duplicate nodes. The second argument
	// of newTrie is the maximal number of branches at any node, which is
	// computed based on the maximal level ml.
	ml := uint8(0)
	for i := range parentIndex {
		l := uint8(parentIndex[i] >> 32)
		if l > ml {
			ml = l
		}
	}
	// One +1 since going one level up; another +1 since counting from zero;
	// the second multiplier is the number of orders on the maximal level.
	trie := newTrie(dc, uint32(ml+1+1)*(1<<uint32(ml)+1))

	cc := uint32(0)

	push := func(p, d uint32, pair uint64) {
		copy(index[cc*dc:], parentIndex[p*dc:(p+1)*dc])
		index[cc*dc+d] = pair

		if !trie.tap(index[cc*dc:]) {
			cc++
		}
	}

	for i := uint32(0); i < pc; i++ {
		for j := uint32(0); j < dc; j++ {
			level := uint8(parentIndex[i*dc+j] >> 32)

			if level == 0 {
				push(i, j, 1<<32|0)
				push(i, j, 1<<32|2)
				continue
			}

			order := uint32(parentIndex[i*dc+j])

			if level == 1 {
				push(i, j, 2<<32|uint64(order+1))
			} else {
				push(i, j, uint64(level+1)<<32|uint64(2*order-1))
				push(i, j, uint64(level+1)<<32|uint64(2*order+1))
			}
		}
	}

	return index[0 : cc*dc]
}
