package newcot

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
		nodes[i] = float64(indices[i]>>32+1) / float64(uint64(2)<<(0xFFFFFFFF&indices[i]))
	}

	return nodes
}

// Refine returns the child indices corresponding to a set of parent indices.
func (o *Open) Refine(indices []uint64) []uint64 {
	nd := o.nd
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
			order := indices[i*nd+j] >> 32

			push(i, j, (level+1)|(2*order)<<32)
			push(i, j, (level+1)|(2*order+2)<<32)
		}
	}

	return childIndices
}

// Balance identifies the missing neighbors of a set of child nodes with respect
// to their parent nodes in each dimension.
func (o *Open) Balance(indices []uint64,
	find func([]uint64) bool, push func([]uint64)) {

	for {
		indices = o.socialize(indices, find, push)
		if len(indices) == 0 {
			break
		}
	}
}

func (o *Open) socialize(indices []uint64,
	find func([]uint64) bool, push func([]uint64)) []uint64 {

	nd := o.nd
	nn := len(indices) / nd

	index := make([]uint64, nd)
	missing := make([]uint64, 0, nd)

	for i := 0; i < nn; i++ {
		copy(index, indices[i*nd:(i+1)*nd])

		for j, pair := range index {
			level := 0xFFFFFFFF & pair
			if level == 0 {
				continue
			}

			level -= 1
			order := (pair >> 32) / 2

			right := order%2 == 1
			if right {
				order = (2*order - 2) / 2
			}

			index[j] = level | order<<32

			if find(index) {
				if right {
					index[j] = (level + 1) | (2*order)<<32
				} else {
					index[j] = (level + 1) | (2*order+2)<<32
				}
				if !find(index) {
					push(index)
					missing = append(missing, index...)
				}
			}

			index[j] = pair
		}
	}

	return missing
}
