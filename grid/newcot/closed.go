package newcot

// Closed represents an instance of the grid in [0, 1]^n.
type Closed struct {
	nd int
}

// NewClosed creates an instance of the grid in [0, 1]^n.
func NewClosed(dimensions uint) *Closed {
	return &Closed{int(dimensions)}
}

// Compute returns the nodes corresponding to the given indices.
func (_ *Closed) Compute(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))

	for i := range nodes {
		level := 0xFFFFFFFF & indices[i]
		if level == 0 {
			nodes[i] = 0.5
		} else {
			nodes[i] = float64(indices[i]>>32) / float64(uint64(2)<<(level-1))
		}
	}

	return nodes
}

// ComputeChildren returns the indices of the child nodes corresponding to the
// parent nodes given by their indices.
func (c *Closed) ComputeChildren(parentIndices []uint64) []uint64 {
	nd := c.nd
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

			if level == 0 {
				push(i, j, 1|0<<32)
				push(i, j, 1|2<<32)
				continue
			}

			order := parentIndices[i*nd+j] >> 32

			if level == 1 {
				push(i, j, 2|(order+1)<<32)
			} else {
				push(i, j, (level+1)|(2*order-1)<<32)
				push(i, j, (level+1)|(2*order+1)<<32)
			}
		}
	}

	return indices[0 : nc*nd]
}
