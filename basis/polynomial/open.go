package polynomial

import (
	"math"

	"github.com/ready-steady/adapt/grid/equidistant"
	"github.com/ready-steady/adapt/internal"
)

// Open is a basis in (0, 1)^n.
type Open struct {
	nd   uint
	grid equidistant.Open
}

// NewOpen creates a basis.
func NewOpen(dimensions, power uint) *Open {
	if power != 1 {
		panic("not implemented")
	}
	return &Open{
		nd:   dimensions,
		grid: *equidistant.NewOpen(1),
	}
}

// Compute evaluates a basis function.
func (self *Open) Compute(index []uint64, point []float64) float64 {
	nd, value := self.nd, 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= self.compute(index[i]&internal.LEVEL_MASK,
			index[i]>>internal.LEVEL_SIZE, point[i])
	}
	return value
}

// Integrate computes the integral of a basis function.
func (self *Open) Integrate(index []uint64) float64 {
	nd, value := self.nd, 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= self.integrate(index[i]&internal.LEVEL_MASK,
			index[i]>>internal.LEVEL_SIZE)
	}
	return value
}

func (self *Open) compute(level, order uint64, x float64) float64 {
	if level == 0 {
		return 1.0
	}

	xi, h := self.grid.Node(level, order)
	count := uint64(2)<<level - 1

	switch order {
	case 0:
		if x >= 2.0*h {
			return 0.0
		}
		return 2.0 - x/h
	case count - 1:
		left := float64(count - 1)
		if x <= left*h {
			return 0.0
		}
		return x/h - left
	default:
		Δ := math.Abs(x - xi)
		if Δ >= h {
			return 0.0
		}
		return 1.0 - Δ/h
	}
}

func (_ *Open) integrate(level, order uint64) float64 {
	if level == 0 {
		return 1.0
	}

	count := uint64(2)<<level - 1

	switch order {
	case 0, count - 1:
		return 2.0 / float64(count+1)
	default:
		return 1.0 / float64(count+1)
	}
}
