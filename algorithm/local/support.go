package local

func assess(basis Basis, target Target, progress *Progress, indices []uint64,
	nodes, surpluses []float64, ni, no uint) []float64 {

	nn := uint(len(indices)) / ni
	scores := measure(basis, indices, ni)
	for i := uint(0); i < nn; i++ {
		location := Location{
			Node:    nodes[i*ni : (i+1)*ni],
			Surplus: surpluses[i*no : (i+1)*no],
			Volume:  scores[i],
		}
		scores[i] = target.Score(&location, progress)
	}

	return scores
}

func cumulate(basis Basis, indices []uint64, surpluses []float64, ni, no uint,
	integral []float64) {

	nn := uint(len(indices)) / ni
	for i := uint(0); i < nn; i++ {
		volume := basis.Integrate(indices[i*ni : (i+1)*ni])
		for j := uint(0); j < no; j++ {
			integral[j] += surpluses[i*no+j] * volume
		}
	}
}

func levelize(indices []uint64, ni uint) []uint {
	nn := uint(len(indices)) / ni
	levels := make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < ni; j++ {
			levels[i] += uint(LEVEL_MASK & indices[i*ni+j])
		}
	}
	return levels
}

func measure(basis Basis, indices []uint64, ni uint) []float64 {
	nn := uint(len(indices)) / ni
	volumes := make([]float64, nn)
	for i := uint(0); i < nn; i++ {
		volumes[i] = basis.Integrate(indices[i*ni : (i+1)*ni])
	}
	return volumes
}
