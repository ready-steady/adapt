package global

func maxFloat64Set(data []float64, set Set) (float64, uint) {
	value, position := -infinity, ^uint(0)
	for i := range set {
		if data[i] > value {
			value, position = data[i], i
		}
	}
	return value, position
}

func maxUint64(data []uint64) uint64 {
	result := uint64(0)
	for _, value := range data {
		if value > result {
			result = value
		}
	}
	return result
}

func minUint64Set(data []uint64, set Set) (uint64, uint) {
	value, position := ^uint64(0), ^uint(0)
	for i := range set {
		if data[i] < value {
			value, position = data[i], i
		}
	}
	return value, position
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
