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
func (self *Open) Compute(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))
	for i := range nodes {
		level := indices[i] & internal.LEVEL_MASK
		order := indices[i] >> internal.LEVEL_SIZE
		nodes[i], _, _ = self.Node(level, order)
	}
	return nodes
}

// Index returns the nodal indices of a set of level indices.
func (self *Open) Index(lindices []uint64) []uint64 {
	return index(lindices, openIndex, self.nd)
}

// Node returns the node corresponding to an index in one dimension.
func (_ *Open) Node(level, order uint64) (node, step float64, count uint64) {
	count = uint64(2)<<level - 1
	step = 1.0 / float64(count+1)
	node = float64(order+1) * step
	return
}

// Parent returns the parent index of an index in one dimension.
func (_ *Open) Parent(level, order uint64) (uint64, uint64) {
	switch level {
	case 0:
		panic("the root does not have a parent")
	default:
		level -= 1
		if order%4 == 0 {
			order = order / 2
		} else {
			order = (order - 2) / 2
		}
	}
	return level, order
}

// Refine returns the child indices of a set of indices.
func (self *Open) Refine(indices []uint64) []uint64 {
	return openRefine(indices, self.nd, 0, self.nd)
}

// RefineToward returns the child indices of a set of indices with respect to a
// particular dimension.
func (self *Open) RefineToward(indices []uint64, i uint) []uint64 {
	return openRefine(indices, self.nd, i, i+1)
}

func openIndex(level uint64) []uint64 {
	if level>>internal.LEVEL_SIZE != 0 {
		panic(fmt.Sprintf("the level %d is too large", level))
	}
	switch level {
	case 0:
		return []uint64{0 | 0<<internal.LEVEL_SIZE}
	default:
		nn := uint(2) << uint(level-1)
		indices := make([]uint64, nn)
		for i := uint(0); i < nn; i++ {
			indices[i] = level | uint64(2*i+1)<<internal.LEVEL_SIZE
		}
		return indices
	}
}

func openRefine(indices []uint64, nd, fd, ld uint) []uint64 {
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
			order := indices[i*nd+j] >> internal.LEVEL_SIZE
			push(i, j, level+1, 2*order)
			push(i, j, level+1, 2*order+2)
		}
	}

	return children[:nc*nd]
}
