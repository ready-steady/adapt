package newton_cotes

func ComputeLevelOrders(level uint8) []uint16 {
	switch level {
	case 0:
		return []uint16{0}
	case 1:
		return []uint16{0, 2}
	}

	count := uint16(2) << (level - 1 - 1)
	orders := make([]uint16, count)

	for i := range orders {
		orders[i] = uint16(2 * (i + 1) - 1)
	}

	return orders
}

func ComputeNodes(levels []uint8, orders []uint16) []float64 {
	count := len(levels)
	nodes := make([]float64, count)

	for i := 0; i < count; i++ {
		if levels[i] > 0 {
			nodes[i] = float64(orders[i]) /
				float64(uint32(2) << (levels[i] - 1))
		} else {
			nodes[i] = 0.5
		}
	}

	return nodes
}
