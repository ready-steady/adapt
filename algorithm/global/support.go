package global

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

func score(basis Basis, strategy strategy, target Target, counts []uint, indices []uint64,
	values, surpluses []float64, ni, no uint) {

	for _, count := range counts {
		oi, oo := count*ni, count*no
		element := Element{
			Indices:   indices[:oi],
			Volumes:   internal.Measure(basis, indices[:oo], ni),
			Values:    values[:oo],
			Surpluses: surpluses[:oo],
		}
		strategy.Push(&element, target.Score(&element))
		indices, values, surpluses = indices[oi:], values[oo:], surpluses[oo:]
	}
}
