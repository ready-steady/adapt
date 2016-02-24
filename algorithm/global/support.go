package global

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

func assess(basis Basis, target Target, counts []uint, indices []uint64,
	values, surpluses []float64, ni, no uint) []float64 {

	scores := make([]float64, len(counts))
	for i, count := range counts {
		offset := count * no
		scores[i] = target.Score(&Location{
			Values:    values[:offset],
			Surpluses: surpluses[:offset],
			Volumes:   internal.Measure(basis, indices[:offset], ni),
		})
		indices, values, surpluses = indices[count:], values[offset:], surpluses[offset:]
	}
	return scores
}
