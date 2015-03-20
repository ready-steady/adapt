package newcot

func compose(levels []uint32, orders []uint32) []uint64 {
	indices := make([]uint64, len(levels))

	for i := range levels {
		indices[i] = uint64(levels[i]) | uint64(orders[i])<<32
	}

	return indices
}
