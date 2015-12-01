package newcot

import (
	"fmt"
)

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

// Children returns the child indices corresponding to a set of parent indices.
func (self *Closed) Children(indices []uint64) []uint64 {
	nd := self.nd
	nn := len(indices) / nd

	childIndices := make([]uint64, 2*nn*nd*nd)

	nc := 0
	push := func(p, d int, level, order uint64) {
		if level>>LEVEL_SIZE != 0 || order>>ORDER_SIZE != 0 {
			panic(fmt.Sprintf("the level %d and order %d overflow uint64", level, order))
		}
		copy(childIndices[nc*nd:], indices[p*nd:(p+1)*nd])
		childIndices[nc*nd+d] = level | order<<LEVEL_SIZE
		nc++
	}

	for i := 0; i < nn; i++ {
		for j := 0; j < nd; j++ {
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

	return childIndices[:nc*nd]
}

// Parent transforms an index into its parent index in the ith dimension.
func (_ *Closed) Parent(index []uint64, i uint) {
	switch level := LEVEL_MASK & index[i]; level {
	case 0:
	case 1:
		index[i] = 0 | 0<<LEVEL_SIZE
	case 2:
		index[i] = 1 | (index[i]>>LEVEL_SIZE-1)<<LEVEL_SIZE
	default:
		order := (index[i]>>LEVEL_SIZE - 1) / 2
		if order%2 == 0 {
			order = (index[i]>>LEVEL_SIZE + 1) / 2
		}
		index[i] = (level - 1) | order<<LEVEL_SIZE
	}
}

// Sibling transforms an index into its sibling index in the ith dimension.
func (_ *Closed) Sibling(index []uint64, i uint) {
	switch level := LEVEL_MASK & index[i]; level {
	case 0, 2:
	case 1:
		if index[i]>>LEVEL_SIZE == 0 {
			index[i] = 1 | 2<<LEVEL_SIZE
		} else {
			index[i] = 1 | 0<<LEVEL_SIZE
		}
	default:
		order := index[i] >> LEVEL_SIZE
		if ((order-1)/2)%2 == 1 {
			order -= 2
		} else {
			order += 2
		}
		index[i] = level | order<<LEVEL_SIZE
	}
}
