package global

func score(target Target, indices []uint64, volumes, observations, surpluses []float64,
	counts []uint, ni, no uint) []float64 {

	nn := uint(len(counts))
	scores := make([]float64, nn)
	for i, offset := uint(0), uint(0); i < nn; i++ {
		fi, fo := offset*ni, offset*no
		li, lo := fi+counts[i]*ni, fo+counts[i]*no
		element := Element{
			Indices:   indices[fi:li],
			Volumes:   volumes[offset:(offset + counts[i])],
			Values:    observations[fo:lo],
			Surpluses: surpluses[fo:lo],
		}
		scores[i] = target.Score(&element)
		offset += counts[i]
	}
	return scores
}
