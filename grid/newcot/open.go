package newcot

// Open represents an instance of the grid on (0, 1)^n.
type Open struct {
	ic uint16
}

// NewOpen creates an instance of the grid on (0, 1)^n.
func NewOpen(dimensions uint16) *Open {
	return &Open{dimensions}
}

// Dimensions returns the number of dimensions of the grid.
func (o *Open) Dimensions() uint16 {
	return o.ic
}

// ComputeNodes returns the nodes corresponding to the given indices.
func (_ *Open) ComputeNodes(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))

	for i := range nodes {
		nodes[i] = float64(indices[i]>>32+1) / float64(uint32(2)<<uint32(indices[i]))
	}

	return nodes
}

// ComputeChildren returns the indices of the child nodes corresponding to the
// parent nodes given by their indices.
func (o *Open) ComputeChildren(parentIndices []uint64) []uint64 {
	ic := uint32(o.ic)
	pc := uint32(len(parentIndices)) / ic

	indices := make([]uint64, 2*pc*ic*ic)

	hash := newHash(ic, 2*pc*ic)

	cc := uint32(0)

	push := func(p, d uint32, pair uint64) {
		copy(indices[cc*ic:], parentIndices[p*ic:(p+1)*ic])
		indices[cc*ic+d] = pair

		if !hash.tap(indices[cc*ic:]) {
			cc++
		}
	}

	var i, j, level, order uint32

	for i = 0; i < pc; i++ {
		for j = 0; j < ic; j++ {
			level = uint32(parentIndices[i*ic+j])
			order = uint32(parentIndices[i*ic+j] >> 32)

			push(i, j, uint64(level+1)|uint64(2*order)<<32)
			push(i, j, uint64(level+1)|uint64(2*order+2)<<32)
		}
	}

	return indices[0 : cc*ic]
}
