package equidistant

func compose(levels []uint64, orders []uint64) []uint64 {
	indices := make([]uint64, len(levels))

	for i := range levels {
		indices[i] = uint64(levels[i]) | uint64(orders[i])<<LEVEL_SIZE
	}

	return indices
}
