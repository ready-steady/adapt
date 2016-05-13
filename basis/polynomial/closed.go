package polynomial

import (
	"math"

	"github.com/ready-steady/adapt/grid/equidistant"
)

// Closed is a basis in [0, 1]^n.
type Closed struct {
	nd   uint
	np   uint
	grid equidistant.Closed
}

// NewClosed creates a basis.
func NewClosed(dimensions uint, power uint) *Closed {
	return &Closed{
		nd:   dimensions,
		np:   power,
		grid: *equidistant.NewClosed(1),
	}
}

// Compute evaluates a basis function.
func (self *Closed) Compute(index []uint64, point []float64) float64 {
	return compute(index, point, self.nd, self.compute)
}

// Integrate computes the integral of a basis function.
func (self *Closed) Integrate(index []uint64) float64 {
	return integrate(index, self.nd, self.integrate)
}

func (self *Closed) compute(level, order uint64, x float64) float64 {
	np := self.np
	if level < uint64(np) {
		np = uint(level)
	}
	if np == 0 {
		return 1.0
	}

	xi, h, _ := self.grid.Node(level, order)

	Δ := math.Abs(x - xi)
	if Δ >= h {
		return 0.0
	}

	if np == 1 {
		// Use two linear segments. The reason is that, taking into account the
		// endpoints, there are three points available in order to construct a
		// first-order polynomial; however, such a polynomial can satisfy only
		// any two of them.
		return 1.0 - Δ/h
	}

	value := 1.0

	// The left endpoint of the local support.
	xl := xi - h
	value *= (x - xl) / (xi - xl)
	np -= 1

	// The right endpoint of the local support.
	xr := xi + h
	value *= (x - xr) / (xi - xr)
	np -= 1

	// Find the rest of the needed ancestors.
	for np > 0 {
		level, order = self.grid.Parent(level, order)
		xj, _, _ := self.grid.Node(level, order)
		if equal(xj, xl) || equal(xj, xr) {
			continue
		}
		value *= (x - xj) / (xi - xj)
		np -= 1
	}

	return value
}

func (self *Closed) integrate(level, order uint64) float64 {
	np := self.np
	if level < uint64(np) {
		np = uint(level)
	}
	if np == 0 {
		return 1.0
	}

	x, h, _ := self.grid.Node(level, order)

	if np == 1 {
		// Use two liner segments. See the corresponding comment in
		// closedCompute.
		if level == 1 {
			return 0.25
		} else {
			return h
		}
	}

	// Use a Gauss–Legendre quadrature rule to integrate. Such a rule with n
	// nodes integrates exactly polynomials up to order 2*n - 1.
	nodes := uint(math.Ceil((float64(np) + 1.0) / 2.0))
	return quadrature(x-h, x+h, nodes, func(x float64) float64 {
		return self.compute(level, order, x)
	})
}
