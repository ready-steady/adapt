package global

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

func index(grid Grid, lindices []uint64, ni uint) ([]uint64, []uint) {
	nn := uint(len(lindices)) / ni
	indices, counts := []uint64(nil), make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		newIndices := grid.Index(lindices[:ni])
		indices = append(indices, newIndices...)
		counts[i] = uint(len(newIndices)) / ni
		lindices = lindices[ni:]
	}
	return indices, counts
}

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
