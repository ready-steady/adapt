package newcot

// Open represents an instance of the grid on (0, 1)^n.
type Open struct {
	dc uint16
}

// NewOpen creates an instance of the grid on (0, 1)^n.
func NewOpen(dimensions uint16) *Open {
	return &Open{dimensions}
}

// Dimensions returns the number of dimensions of the grid.
func (o *Open) Dimensions() uint16 {
	return o.dc
}

// ComputeNodes returns the nodes corresponding to the given index.
func (_ *Open) ComputeNodes(index []uint64) []float64 {
	nodes := make([]float64, len(index))

	for i := range nodes {
		nodes[i] = float64(index[i]>>32+1) / float64(uint32(2)<<uint32(index[i]))
	}

	return nodes
}

// ComputeChildren returns the index of the child nodes corresponding to the
// parent nodes given by their index.
func (o *Open) ComputeChildren(parentIndex []uint64) []uint64 {
	dc := uint32(o.dc)
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
			order = uint32(parentIndex[i*dc+j] >> 32)

			push(i, j, uint64(level+1)|uint64(2*order)<<32)
			push(i, j, uint64(level+1)|uint64(2*order+2)<<32)
		}
	}

	return index[0 : cc*dc]
}
