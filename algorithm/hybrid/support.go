package hybrid

func assess(basis Basis, target Target, counts []uint, indices []uint64,
	values, surpluses []float64, ni, no uint) []float64 {

	scores := make([]float64, len(counts))
	for i, count := range counts {
		offset := count * no
		scores[i] = target.Score(&Location{
			Values:    values[:offset],
			Surpluses: surpluses[:offset],
			Volumes:   measure(basis, indices[:offset], ni),
		})
		indices, values, surpluses = indices[count:], values[offset:], surpluses[offset:]
	}
	return scores
}

func measure(basis Basis, indices []uint64, ni uint) []float64 {
	nn := uint(len(indices)) / ni
	volumes := make([]float64, nn)
	for i := uint(0); i < nn; i++ {
		volumes[i] = basis.Integrate(indices[i*ni : (i+1)*ni])
	}
	return volumes
}
