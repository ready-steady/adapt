package local

import (
	"github.com/ready-steady/adapt/internal"
)

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
			levels[i] += uint(internal.LEVEL_MASK & indices[i*ni+j])
		}
	}
	return levels
}

func score(target Target, state *state, ni, no uint) []float64 {
	nn := uint(len(state.Indices)) / ni
	scores := make([]float64, nn)
	for i := uint(0); i < nn; i++ {
		scores[i] = target.Score(&Element{
			Index:   state.Indices[i*ni : (i+1)*ni],
			Volume:  state.Volumes[i],
			Value:   state.Observations[i*no : (i+1)*no],
			Surplus: state.Surpluses[i*no : (i+1)*no],
		})
	}
	return scores
}
