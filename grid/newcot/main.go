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

	cc := uint32(0)

	push := func(p, d uint32, pair uint64) {
		copy(index[cc*dc:], parentIndex[p*dc:(p+1)*dc])
		index[cc*dc+d] = pair

		// Check uniqueness
		for i := uint32(0); i < cc; i++ {
			found := true

			for j := uint32(0); j < dc; j++ {
				if index[i*dc+j] != index[cc*dc+j] {
					found = false
					break
				}
			}

			if found {
				// Discard since a duplicate
				return
			}
		}

		cc++
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
