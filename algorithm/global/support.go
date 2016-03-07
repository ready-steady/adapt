package global

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

func score(basis Basis, strategy strategy, target Target, lindices, indices []uint64,
	counts []uint, values, surpluses []float64, ni, no uint) {

	for _, count := range counts {
		oi, oo := count*ni, count*no
		element := Element{
			Lindex:    lindices[:ni],
			Indices:   indices[:oi],
			Volumes:   internal.Measure(basis, indices[:oo], ni),
			Values:    values[:oo],
			Surpluses: surpluses[:oo],
		}
		strategy.Push(&element, target.Score(&element))
		lindices, indices = lindices[ni:], indices[oi:]
		values, surpluses = values[oo:], surpluses[oo:]
	}
}
