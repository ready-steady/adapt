package newcot

// Closed represents an instance of the grid on [0, 1]^n.
type Closed struct {
	dc uint
}

// NewClosed creates an instance of the grid on [0, 1]^n.
func NewClosed(dimensions uint) *Closed {
	return &Closed{dimensions}
}

// ComputeNodes returns the nodes corresponding to the given indices.
func (_ *Closed) ComputeNodes(indices []uint64) []float64 {
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
	dc := c.dc
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

			if level == 0 {
				push(i, j, 1|0<<32)
				push(i, j, 1|2<<32)
				continue
			}

			order := parentIndices[i*dc+j] >> 32

			if level == 1 {
				push(i, j, 2|(order+1)<<32)
			} else {
				push(i, j, (level+1)|(2*order-1)<<32)
				push(i, j, (level+1)|(2*order+1)<<32)
			}
		}
	}

	return indices[0 : cc*dc]
}
