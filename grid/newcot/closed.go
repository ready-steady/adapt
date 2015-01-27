package newcot

// Closed represents an instance of the grid on [0, 1]^n.
type Closed struct {
	ic uint16
}

// NewClosed creates an instance of the grid on [0, 1]^n.
func NewClosed(dimensions uint16) *Closed {
	return &Closed{dimensions}
}

// Dimensions returns the number of dimensions of the grid.
func (c *Closed) Dimensions() uint16 {
	return c.ic
}

// ComputeNodes returns the nodes corresponding to the given indices.
func (_ *Closed) ComputeNodes(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))

	for i := range nodes {
		if uint32(indices[i]) == 0 {
			nodes[i] = 0.5
		} else {
			nodes[i] = float64(indices[i]>>32) / float64(uint32(2)<<(uint32(indices[i])-1))
		}
	}

	return nodes
}

// ComputeChildren returns the indices of the child nodes corresponding to the
// parent nodes given by their indices.
func (c *Closed) ComputeChildren(parentIndices []uint64) []uint64 {
	ic := uint32(c.ic)
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

			if level == 0 {
				push(i, j, 1|0<<32)
				push(i, j, 1|2<<32)
				continue
			}

			order = uint32(parentIndices[i*ic+j] >> 32)

			if level == 1 {
				push(i, j, 2|uint64(order+1)<<32)
			} else {
				push(i, j, uint64(level+1)|uint64(2*order-1)<<32)
				push(i, j, uint64(level+1)|uint64(2*order+1)<<32)
			}
		}
	}

	return indices[0 : cc*ic]
}
