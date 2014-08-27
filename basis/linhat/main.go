// Package linhat provides means for working with the basis formed by the
// linear hat function.
package linhat

import (
	"math"
)

// Self represents a particular instantiation of the basis.
type Self struct {
}

// New creates an instance of the basis.
func New() *Self {
	return &Self{}
}

// Evaluate computes the value of the basis function, identified by the given
// level and order, at the given point.
func (_ *Self) Evaluate(point float64, level uint8, order uint32) float64 {
	if level == 0 {
		if math.Abs(point-0.5) > 0.5 {
			return 0
		} else {
			return 1
		}
	}

	limit := float64(uint32(2) << (level - 1))
	delta := math.Abs(point - float64(order)/limit)

	if delta >= 1/limit {
		return 0
	} else {
		return 1 - limit*delta
	}
}
