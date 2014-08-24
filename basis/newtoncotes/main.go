package newton_cotes

func ComputeOrders(level uint8) []uint16 {
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

func ComputeChildren(parentLevels []uint8, parentOrders []uint16) ([]uint8, []uint16) {
	parentCount := len(parentLevels)

	levels := make([]uint8, 2 * parentCount)
	orders := make([]uint16, 2 * parentCount)

	k := 0

	for i := 0; i < parentCount; i++ {
		switch parentLevels[i] {
		case 0:
			levels[k] = 1
			orders[k] = 0
			k++

			levels[k] = 1
			orders[k] = 2
			k++

		case 1:
			levels[k] = 2
			orders[k] = parentOrders[i] + 1
			k++

		default:
			levels[k] = parentLevels[i] + 1
			orders[k] = 2 * parentOrders[i] - 1
			k++

			levels[k] = parentLevels[i] + 1
			orders[k] = 2 * parentOrders[i] + 1
			k++
		}
	}

	return levels[0:k], orders[0:k]
}
