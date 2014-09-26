package newcot

func compose(levels []uint32, orders []uint32) []uint64 {
	index := make([]uint64, len(levels))

	for i := range levels {
		index[i] = uint64(levels[i]) | uint64(orders[i])<<32
	}

	return index
}
