package global

func find(data []bool) []uint {
	n := uint(len(data))
	indices := make([]uint, 0, n)
	for i := uint(0); i < n; i++ {
		if data[i] {
			indices = append(indices, i)
		}
	}
	return indices
}

func maxFloat64(data []float64, indices ...uint) (value float64, position uint) {
	if indices == nil {
		count := uint(len(data))
		value, position = data[0], 0
		for i := uint(1); i < count; i++ {
			if data[i] > value {
				value, position = data[i], i
			}
		}
	} else {
		count := uint(len(indices))
		value, position = data[indices[0]], indices[0]
		for i := uint(1); i < count; i++ {
			j := indices[i]
			if data[j] > value {
				value, position = data[j], j
			}
		}
	}
	return
}

func maxUint(data []uint, indices ...uint) (value uint, position uint) {
	if indices == nil {
		count := uint(len(data))
		value, position = data[0], 0
		for i := uint(1); i < count; i++ {
			if data[i] > value {
				value, position = data[i], i
			}
		}
	} else {
		count := uint(len(indices))
		value, position = data[indices[0]], indices[0]
		for i := uint(1); i < count; i++ {
			j := indices[i]
			if data[j] > value {
				value, position = data[j], j
			}
		}
	}
	return
}

func minUint(data []uint, indices ...uint) (value uint, position uint) {
	if indices == nil {
		count := uint(len(data))
		value, position = data[0], 0
		for i := uint(1); i < count; i++ {
			if data[i] < value {
				value, position = data[i], i
			}
		}
	} else {
		count := uint(len(indices))
		value, position = data[indices[0]], indices[0]
		for i := uint(1); i < count; i++ {
			j := indices[i]
			if data[j] < value {
				value, position = data[j], j
			}
		}
	}
	return
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
