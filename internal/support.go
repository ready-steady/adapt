package internal

// Compose encodes levels and orders.
func Compose(levels []uint64, orders []uint64) (indices []uint64) {
	indices = make([]uint64, len(levels))
	for i := range levels {
		indices[i] = uint64(levels[i]) | uint64(orders[i])<<LEVEL_SIZE
	}
	return
}

// Decompose decodes levels and orders.
func Decompose(indices []uint64) (levels []uint64, orders []uint64) {
	levels, orders = make([]uint64, len(indices)), make([]uint64, len(indices))
	for i := range indices {
		levels[i], orders[i] = indices[i]&LEVEL_MASK, indices[i]>>LEVEL_SIZE
	}
	return
}
