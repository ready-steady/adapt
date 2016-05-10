package equidistant

import (
	"fmt"

	"github.com/ready-steady/adapt/internal"
)

// Closed is a grid in [0, 1]^n.
type Closed struct {
	nd uint
}

// NewClosed creates a grid.
func NewClosed(dimensions uint) *Closed {
	return &Closed{dimensions}
}

// Compute returns the nodes corresponding to a set of indices.
func (self *Closed) Compute(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))
	for i := range nodes {
		level := indices[i] & internal.LEVEL_MASK
		order := indices[i] >> internal.LEVEL_SIZE
		nodes[i], _ = self.Node(level, order)
	}
	return nodes
}

// Index returns the nodal indices of a set of level indices.
func (self *Closed) Index(lindices []uint64) []uint64 {
	return index(lindices, closedIndex, self.nd)
}

// Node returns the node corresponding to an index in one dimension.
func (_ *Closed) Node(level, order uint64) (node, step float64) {
	if level == 0 {
		node = 0.5
		step = 0.5
	} else {
		step = 1.0 / float64(uint64(2)<<(level-1))
		node = float64(order) * step
	}
	return
}

// Parent returns the parent index of an index in one dimension.
func (_ *Closed) Parent(level, order uint64) (uint64, uint64) {
	switch level {
	case 0:
		panic("the root does not have a parent")
	case 1:
		level = 0
		order = 0
	case 2:
		level = 1
		order -= 1
	default:
		level -= 1
		if (order-1)%4 == 0 {
			order = (order + 1) / 2
		} else {
			order = (order - 1) / 2
		}
	}
	return level, order
}

// Refine returns the child indices of a set of indices.
func (self *Closed) Refine(indices []uint64) []uint64 {
	return closedRefine(indices, self.nd, 0, self.nd)
}

// RefineToward returns the child indices of a set of indices with respect to a
// particular dimension.
func (self *Closed) RefineToward(indices []uint64, i uint) []uint64 {
	return closedRefine(indices, self.nd, i, i+1)
}

func closedIndex(level uint64) []uint64 {
	if level>>internal.LEVEL_SIZE != 0 {
		panic(fmt.Sprintf("the level %d is too large", level))
	}
	switch level {
	case 0:
		return []uint64{0 | 0<<internal.LEVEL_SIZE}
	case 1:
		return []uint64{1 | 0<<internal.LEVEL_SIZE, 1 | 2<<internal.LEVEL_SIZE}
	default:
		nn := uint(2) << uint(level-2)
		indices := make([]uint64, nn)
		for i := uint(0); i < nn; i++ {
			indices[i] = level | uint64(2*i+1)<<internal.LEVEL_SIZE
		}
		return indices
	}
}

func closedRefine(indices []uint64, nd, fd, ld uint) []uint64 {
	nn := uint(len(indices)) / nd

	children := make([]uint64, 2*nn*nd*(ld-fd))

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
		for j := fd; j < ld; j++ {
			level := indices[i*nd+j] & internal.LEVEL_MASK

			if level == 0 {
				push(i, j, 1, 0)
				push(i, j, 1, 2)
				continue
			}

			order := indices[i*nd+j] >> internal.LEVEL_SIZE

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
