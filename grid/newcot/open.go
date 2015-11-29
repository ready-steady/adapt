package newcot

import (
	"fmt"
)

// Open represents an instance of the grid in (0, 1)^n.
type Open struct {
	nd int
}

// NewOpen creates an instance of the grid in (0, 1)^n.
func NewOpen(dimensions uint) *Open {
	return &Open{int(dimensions)}
}

// Compute returns the nodes corresponding to the given indices.
func (_ *Open) Compute(indices []uint64) []float64 {
	nodes := make([]float64, len(indices))
	for i := range nodes {
		nodes[i] = float64(indices[i]>>LEVEL_SIZE+1) / float64(uint64(2)<<(LEVEL_MASK&indices[i]))
	}
	return nodes
}

// Refine returns the child indices corresponding to a set of parent indices.
func (self *Open) Refine(indices []uint64) []uint64 {
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
			order := indices[i*nd+j] >> LEVEL_SIZE
			push(i, j, level+1, 2*order)
			push(i, j, level+1, 2*order+2)
		}
	}

	return childIndices
}

// Parent transforms an index into its parent index in the ith dimension.
func (_ *Open) Parent(index []uint64, i uint) {
	level := LEVEL_MASK & index[i]
	if level == 0 {
		return
	}

	order := (index[i] >> LEVEL_SIZE) / 2
	if order%2 == 1 {
		order = (index[i]>>LEVEL_SIZE - 2) / 2
	}

	index[i] = (level - 1) | order<<LEVEL_SIZE
}

// Sibling transforms an index into its sibling index in the ith dimension.
func (_ *Open) Sibling(index []uint64, i uint) {
	level := LEVEL_MASK & index[i]
	if level == 0 {
		return
	}

	order := index[i] >> LEVEL_SIZE
	if (order/2)%2 == 1 {
		order -= 2
	} else {
		order += 2
	}

	index[i] = level | order<<LEVEL_SIZE
}
