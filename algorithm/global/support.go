package global

func assess(target Target, progress *Progress, surpluses []float64,
	counts []uint, no uint) []float64 {

	scores := make([]float64, 0, len(counts))
	for _, count := range counts {
		location := Location{
			Surpluses: surpluses[:count*no],
		}
		scores = append(scores, target.Score(&location, progress))
		surpluses = surpluses[count*no:]
	}
	return scores
}

func maxFloat64(data []float64, cursor cursor) (float64, uint) {
	value, position := -infinity, ^uint(0)
	for i := range cursor {
		if data[i] > value {
			value, position = data[i], i
		}
	}
	return value, position
}

func maxUint(data []uint) (uint, uint) {
	count := uint(len(data))
	value, position := data[0], uint(0)
	for i := uint(1); i < count; i++ {
		if data[i] > value {
			value, position = data[i], i
		}
	}
	return value, position
}

func minUint(data []uint, cursor cursor) (uint, uint) {
	value, position := ^uint(0), ^uint(0)
	for i := range cursor {
		if data[i] < value {
			value, position = data[i], i
		}
	}
	return value, position
}

func repeatBool(value bool, times uint) []bool {
	data := make([]bool, times)
	for i := uint(0); i < times; i++ {
		data[i] = value
	}
	return data
}

func repeatFloat64(value float64, times uint) []float64 {
	data := make([]float64, times)
	for i := uint(0); i < times; i++ {
		data[i] = value
	}
	return data
}

func repeatUint(value uint, times uint) []uint {
	data := make([]uint, times)
	for i := uint(0); i < times; i++ {
		data[i] = value
	}
	return data
}

func repeatUint8(value uint8, times uint) []uint8 {
	data := make([]uint8, times)
	for i := uint(0); i < times; i++ {
		data[i] = value
	}
	return data
}
