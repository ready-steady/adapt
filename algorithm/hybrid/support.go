package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

func score(basis Basis, target Target, indices []uint64, counts []uint,
	values, surpluses []float64, ni, no uint) []float64 {

	nn := uint(len(counts))
	scores := make([]float64, 0, no)
	for i, offset := uint(0), uint(0); i < nn; i++ {
		fi, fo := offset*ni, offset*no
		li, lo := fi+counts[i]*ni, fo+counts[i]*no
		element := Element{
			Indices:   indices[fi:li],
			Volumes:   internal.Measure(basis, indices[fo:lo], ni),
			Values:    values[fo:lo],
			Surpluses: surpluses[fo:lo],
		}
		scores = append(scores, target.Score(&element)...)
		offset += counts[i]
	}
	return scores
}
