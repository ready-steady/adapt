package newcot

// Open represents an instance of the grid on (0, 1)^n.
type Open struct {
	dc uint
}

// NewOpen creates an instance of the grid on (0, 1)^n.
func NewOpen(dimensions uint) *Open {
	return &Open{dimensions}
}

// ComputeNodes returns the nodes corresponding to the given indices.
func (_ *Open) ComputeNodes(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))

	for i := range nodes {
		nodes[i] = float64(indices[i]>>32+1) / float64(uint64(2)<<(0xFFFFFFFF&indices[i]))
	}

	return nodes
}

// ComputeChildren returns the indices of the child nodes corresponding to the
// parent nodes given by their indices.
func (o *Open) ComputeChildren(parentIndices []uint64) []uint64 {
	dc := o.dc
	pc := uint(len(parentIndices)) / dc

	indices := make([]uint64, 2*pc*dc*dc)

	hash := newHash(dc, 2*pc*dc)

	cc := uint(0)

	push := func(p, d uint, pair uint64) {
		copy(indices[cc*dc:], parentIndices[p*dc:(p+1)*dc])
		indices[cc*dc+d] = pair

		if !hash.tap(indices[cc*dc:]) {
			cc++
		}
	}

	for i := uint(0); i < pc; i++ {
		for j := uint(0); j < dc; j++ {
			level := 0xFFFFFFFF & parentIndices[i*dc+j]
			order := parentIndices[i*dc+j] >> 32

			push(i, j, (level+1)|(2*order)<<32)
			push(i, j, (level+1)|(2*order+2)<<32)
		}
	}

	return indices[0 : cc*dc]
}
