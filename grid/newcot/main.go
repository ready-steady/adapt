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

// ComputeNodes returns the nodes corresponding to the given levels and orders.
func (_ *Self) ComputeNodes(levels []uint8, orders []uint32) []float64 {
	nodes := make([]float64, len(levels))

	for i := range nodes {
		if levels[i] == 0 {
			nodes[i] = 0.5
		} else {
			nodes[i] = float64(orders[i]) / float64(uint32(2)<<(levels[i]-1))
		}
	}

	return nodes
}

// ComputeChildren returns the levels and orders of the child nodes
// corresponding to the parent nodes given by their levels and orders.
func (self *Self) ComputeChildren(parentLevels []uint8, parentOrders []uint32) ([]uint8, []uint32) {
	dc := uint32(self.dc)
	pc := uint32(len(parentLevels)) / dc

	levels := make([]uint8, 2*pc*dc*dc)
	orders := make([]uint32, 2*pc*dc*dc)

	cc := uint32(0)

	push := func(p, d uint32, level uint8, order uint32) {
		for i := uint32(0); i < dc; i++ {
			if i == d {
				levels[cc*dc+i] = level
				orders[cc*dc+i] = order
			} else {
				levels[cc*dc+i] = parentLevels[p*dc+i]
				orders[cc*dc+i] = parentOrders[p*dc+i]
			}
		}

		// Check uniqueness
		for i := uint32(0); i < cc; i++ {
			found := true

			for j := uint32(0); j < dc; j++ {
				if levels[i*dc+j] != levels[cc*dc+j] ||
					orders[i*dc+j] != orders[cc*dc+j] {

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
			level := parentLevels[i*dc+j]
			order := parentOrders[i*dc+j]

			switch level {
			case 0:
				push(i, j, 1, 0)
				push(i, j, 1, 2)
			case 1:
				push(i, j, 2, order+1)
			default:
				push(i, j, level+1, 2*order-1)
				push(i, j, level+1, 2*order+1)
			}
		}
	}

	return levels[0 : cc*dc], orders[0 : cc*dc]
}
