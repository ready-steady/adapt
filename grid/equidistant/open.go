package equidistant

import (
	"fmt"

	"github.com/ready-steady/adapt/internal"
)

// Open is a grid in (0, 1)^n.
type Open struct {
	nd uint
}

// NewOpen creates a grid.
func NewOpen(dimensions uint) *Open {
	return &Open{dimensions}
}

// Compute returns the nodes corresponding to a set of indices.
func (_ *Open) Compute(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))
	for i := range nodes {
		level := internal.LEVEL_MASK & indices[i]
		order := indices[i] >> internal.LEVEL_SIZE
		nodes[i] = float64(order+1) / float64(uint64(2)<<level)
	}
	return nodes
}

// Refine returns the child indices of a set of indices.
func (self *Open) Refine(indices []uint64) []uint64 {
	nd := self.nd
	nn := uint(len(indices)) / nd

	children := make([]uint64, 2*nn*nd*nd)

	nc := uint(0)
	push := func(p, d uint, level, order uint64) {
		if level>>internal.LEVEL_SIZE != 0 || order>>internal.ORDER_SIZE != 0 {
			panic(fmt.Sprintf("the level %d or order %d is too large", level, order))
		}
		copy(children[nc*nd:], indices[p*nd:(p+1)*nd])
		children[nc*nd+d] = level | order<<internal.LEVEL_SIZE
		nc++
	}

	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nd; j++ {
			level := internal.LEVEL_MASK & indices[i*nd+j]
			order := indices[i*nd+j] >> internal.LEVEL_SIZE
			push(i, j, level+1, 2*order)
			push(i, j, level+1, 2*order+2)
		}
	}

	return children
}
