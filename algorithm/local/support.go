package local

func assess(basis Basis, target Target, indices []uint64,
	values, surpluses []float64, ni, no uint) []float64 {

	nn := uint(len(indices)) / ni
	scores := measure(basis, indices, ni)
	for i := uint(0); i < nn; i++ {
		scores[i] = target.Score(&Location{
			Value:   values[i*no : (i+1)*no],
			Surplus: surpluses[i*no : (i+1)*no],
			Volume:  scores[i],
		})
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

func filter(indices []uint64, scores []float64, lmin, lmax, ni uint) []uint64 {
	nn := uint(len(scores))
	levels := levelize(indices, ni)

	na, ne := uint(0), nn
	for i, j := uint(0), uint(0); i < nn; i++ {
		if levels[i] >= lmin && (scores[i] <= 0.0 || levels[i] >= lmax) {
			j++
			continue
		}
		if j > na {
			copy(indices[na*ni:], indices[j*ni:ne*ni])
			ne -= j - na
			j = na
		}
		na++
		j++
	}

	return indices[:na*ni]
}

func levelize(indices []uint64, ni uint) []uint {
	nn := uint(len(indices)) / ni
	levels := make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < ni; j++ {
			levels[i] += uint(levelMask & indices[i*ni+j])
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
