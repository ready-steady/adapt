package linhat

func compose(level, order uint64) uint64 {
	return level | order<<LEVEL_SIZE
}
