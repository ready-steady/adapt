package polynomial

import (
	"math"

	"github.com/ready-steady/adapt/grid/equidistant"
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
	return compute(index, point, self.nd, self.compute)
}

// Integrate computes the integral of a basis function.
func (self *Open) Integrate(index []uint64) float64 {
	return integrate(index, self.nd, self.integrate)
}

func (self *Open) compute(level, order uint64, x float64) float64 {
	if level == 0 {
		return 1.0
	}
	xi, h, count := self.grid.Node(level, order)
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

func (self *Open) integrate(level, order uint64) float64 {
	if level == 0 {
		return 1.0
	}
	_, _, count := self.grid.Node(level, order)
	switch order {
	case 0, count - 1:
		return 2.0 / float64(count+1)
	default:
		return 1.0 / float64(count+1)
	}
}
