package linhat

func compose(level, order uint32) uint64 {
	return uint64(level) | uint64(order)<<32
}
