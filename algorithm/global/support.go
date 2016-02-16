package global

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
