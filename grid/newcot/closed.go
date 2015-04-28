package newcot

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
		level := 0xFFFFFFFF & indices[i]
		if level == 0 {
			nodes[i] = 0.5
		} else {
			nodes[i] = float64(indices[i]>>32) / float64(uint64(2)<<(level-1))
		}
	}

	return nodes
}

// Refine returns the child indices corresponding to a set of parent indices.
func (c *Closed) Refine(indices []uint64) []uint64 {
	nd := c.nd
	nn := len(indices) / nd

	childIndices := make([]uint64, 2*nn*nd*nd)

	nc := 0
	push := func(p, d int, pair uint64) {
		copy(childIndices[nc*nd:], indices[p*nd:(p+1)*nd])
		childIndices[nc*nd+d] = pair
		nc++
	}

	for i := 0; i < nn; i++ {
		for j := 0; j < nd; j++ {
			level := 0xFFFFFFFF & indices[i*nd+j]

			if level == 0 {
				push(i, j, 1|0<<32)
				push(i, j, 1|2<<32)
				continue
			}

			order := indices[i*nd+j] >> 32

			if level == 1 {
				push(i, j, 2|(order+1)<<32)
			} else {
				push(i, j, (level+1)|(2*order-1)<<32)
				push(i, j, (level+1)|(2*order+1)<<32)
			}
		}
	}

	return childIndices[:nc*nd]
}

// Parent transforms an index into its parent index in the ith dimension.
func (_ *Closed) Parent(index []uint64, i uint) {
	level := 0xFFFFFFFF & index[i]
	if level == 0 {
		return
	}
	level -= 1

	var order uint64
	switch level {
	case 0:
		order = 0
	case 1:
		order = index[i] >> 32
		order -= 1
	default:
		order = (index[i]>>32 - 1) / 2
		if order%2 == 0 {
			order = (index[i]>>32 + 1) / 2
		}
	}

	index[i] = level | order<<32
}

// Sibling transforms an index into its sibling index in the ith dimension.
func (_ *Closed) Sibling(_ []uint64, _ uint) {
}
