package equidistant

import (
	"fmt"
)

// Closed represents an instance of the grid in [0, 1]^n.
type Closed struct {
	nd uint
}

// NewClosed creates an instance of the grid in [0, 1]^n.
func NewClosed(dimensions uint) *Closed {
	return &Closed{dimensions}
}

// Compute returns the nodes corresponding to a set of indices.
func (_ *Closed) Compute(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))
	for i := range nodes {
		level := LEVEL_MASK & indices[i]
		if level == 0 {
			nodes[i] = 0.5
		} else {
			order := indices[i] >> LEVEL_SIZE
			nodes[i] = float64(order) / float64(uint64(2)<<(level-1))
		}
	}
	return nodes
}

// Children returns the child indices corresponding to a set of indices.
func (self *Closed) Children(indices []uint64) []uint64 {
	nd := self.nd
	nn := uint(len(indices)) / nd

	children := make([]uint64, 2*nn*nd*nd)

	nc := uint(0)
	push := func(p, d uint, level, order uint64) {
		if level>>LEVEL_SIZE != 0 || order>>ORDER_SIZE != 0 {
			panic(fmt.Sprintf("the level %d and order %d are too large", level, order))
		}
		copy(children[nc*nd:], indices[p*nd:(p+1)*nd])
		children[nc*nd+d] = level | order<<LEVEL_SIZE
		nc++
	}

	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nd; j++ {
			level := LEVEL_MASK & indices[i*nd+j]

			if level == 0 {
				push(i, j, 1, 0)
				push(i, j, 1, 2)
				continue
			}

			order := indices[i*nd+j] >> LEVEL_SIZE

			if level == 1 {
				push(i, j, 2, order+1)
			} else {
				push(i, j, level+1, 2*order-1)
				push(i, j, level+1, 2*order+1)
			}
		}
	}

	return children[:nc*nd]
}
