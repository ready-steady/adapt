package equidistant

import (
	"fmt"
)

// Closed is a grid in [0, 1]^n.
type Closed struct {
	nd uint
}

// NewClosed creates a grid in [0, 1]^n.
func NewClosed(dimensions uint) *Closed {
	return &Closed{dimensions}
}

// Compute returns the nodes corresponding to a set of indices.
func (_ *Closed) Compute(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))
	for i := range nodes {
		level := levelMask & indices[i]
		if level == 0 {
			nodes[i] = 0.5
		} else {
			order := indices[i] >> levelSize
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
		if level>>levelSize != 0 || order>>orderSize != 0 {
			panic(fmt.Sprintf("the level %d or order %d is too large", level, order))
		}
		copy(children[nc*nd:], indices[p*nd:(p+1)*nd])
		children[nc*nd+d] = level | order<<levelSize
		nc++
	}

	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nd; j++ {
			level := levelMask & indices[i*nd+j]

			if level == 0 {
				push(i, j, 1, 0)
				push(i, j, 1, 2)
				continue
			}

			order := indices[i*nd+j] >> levelSize

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

// Index returns the indices of a set of levels.
func (self *Closed) Index(levels []uint8) []uint64 {
	return index(levels, indexClosed, self.nd)
}

func indexClosed(level uint8) []uint64 {
	if level>>levelMask != 0 {
		panic(fmt.Sprintf("the level %d is too large", level))
	}
	switch level {
	case 0:
		return []uint64{0 | 0<<levelSize}
	case 1:
		return []uint64{1 | 0<<levelSize, 1 | 2<<levelSize}
	default:
		nn := uint(2) << uint(level-2)
		indices := make([]uint64, nn)
		for i := uint(0); i < nn; i++ {
			indices[i] = uint64(level) | uint64(2*i+1)<<levelSize
		}
		return indices
	}
}
