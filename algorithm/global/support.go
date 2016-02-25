package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1.0)
)

func assess(basis Basis, target Target, counts []uint, indices []uint64,
	values, surpluses []float64, ni, no uint) {

	for _, count := range counts {
		oi, oo := count*ni, count*no
		target.Score(&Location{
			Indices:   indices[:oi],
			Volumes:   internal.Measure(basis, indices[:oo], ni),
			Values:    values[:oo],
			Surpluses: surpluses[:oo],
		})
		indices, values, surpluses = indices[oi:], values[oo:], surpluses[oo:]
	}
}
