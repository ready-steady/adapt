// Package newcot provides means for working with the Newton–Cotes grid on a
// unit hypercube including boundaries.
//
// Each node in the grid is identified by a sequence of levels and orders.
// Throughout the package, such a sequence is encoded as a sequence of uint64s,
// referred to as an index, where each uint64 is (level|order<<32).
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

// ComputeNodes returns the nodes corresponding to the given index.
func (_ *Self) ComputeNodes(index []uint64) []float64 {
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
func (self *Self) ComputeChildren(parentIndex []uint64) []uint64 {
	dc := uint32(self.dc)
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
