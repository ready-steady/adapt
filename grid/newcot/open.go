package newcot

// Open represents an instance of the grid in (0, 1)^n.
type Open struct {
	nd int
}

// NewOpen creates an instance of the grid in (0, 1)^n.
func NewOpen(dimensions uint) *Open {
	return &Open{int(dimensions)}
}

// Compute returns the nodes corresponding to the given indices.
func (_ *Open) Compute(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))

	for i := range nodes {
		nodes[i] = float64(indices[i]>>32+1) / float64(uint64(2)<<(0xFFFFFFFF&indices[i]))
	}

	return nodes
}

// ComputeChildren returns the indices of the child nodes corresponding to the
// parent nodes given by their indices.
func (o *Open) ComputeChildren(parentIndices []uint64) []uint64 {
	nd := o.nd
	np := len(parentIndices) / nd

	indices := make([]uint64, 2*np*nd*nd)

	hash := newHash(nd, 2*np*nd)

	nc := 0

	push := func(p, d int, pair uint64) {
		copy(indices[nc*nd:], parentIndices[p*nd:(p+1)*nd])
		indices[nc*nd+d] = pair

		if !hash.tap(indices[nc*nd:]) {
			nc++
		}
	}

	for i := 0; i < np; i++ {
		for j := 0; j < nd; j++ {
			level := 0xFFFFFFFF & parentIndices[i*nd+j]
			order := parentIndices[i*nd+j] >> 32

			push(i, j, (level+1)|(2*order)<<32)
			push(i, j, (level+1)|(2*order+2)<<32)
		}
	}

	return indices[0 : nc*nd]
}
