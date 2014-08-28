package adhier

func makeGrid1D(size uint32) []float64 {
	step := 1 / float64(size-1)
	points := make([]float64, size)
	for i := uint32(0); i < size; i++ {
		points[i] = step * float64(i)
	}
	return points
}

func makeGrid2D(size uint32) []float64 {
	step := 1 / float64(size-1)
	points := make([]float64, 2*size*size)
	for k, i := uint32(0), uint32(0); i < size; i++ {
		for j := uint32(0); j < size; j++ {
			points[k] = step * float64(i)
			k++
			points[k] = step * float64(j)
			k++
		}
	}
	return points
}
