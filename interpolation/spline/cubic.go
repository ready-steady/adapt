package spline

import (
	"github.com/ready-steady/linear/system"
)

// Cubic represents a cubic interpolant.
type Cubic struct {
	breaks  []float64
	weights []float64
}

// NewCubic constructs a cubic interpolant for a series of points. The abscissae
// are assumed to a strictly increasing sequence with at lest two elements.
func NewCubic(x, y []float64) *Cubic {
	n := uint(len(x))

	dx := make([]float64, n-1)
	dydx := make([]float64, n-1)
	for i := uint(0); i < (n - 1); i++ {
		dx[i] = x[i+1] - x[i]
		dydx[i] = (y[i+1] - y[i]) / dx[i]
	}

	s := &Cubic{}

	switch n {
	case 2:
		s.breaks = []float64{x[0], x[1]}
		s.weights = []float64{0, 0, dydx[0], y[0]}
	case 3:
		c1 := (dydx[1] - dydx[0]) / (x[2] - x[0])
		s.breaks = []float64{x[0], x[2]}
		s.weights = []float64{0, c1, dydx[0] - c1*dx[0], y[0]}
	default:
		xb := x[2] - x[0]
		xe := x[n-1] - x[n-3]

		a := make([]float64, n)
		for i := uint(0); i < (n - 2); i++ {
			a[i] = dx[i+1]
		}
		a[n-2] = xe

		b := make([]float64, n)
		b[0] = dx[1]
		for i := uint(1); i < (n - 1); i++ {
			b[i] = 2 * (dx[i] + dx[i-1])
		}
		b[n-1] = dx[n-3]

		c := make([]float64, n)
		c[1] = xb
		for i := uint(2); i < n; i++ {
			c[i] = dx[i-2]
		}

		d := make([]float64, n)
		d[0] = ((dx[0]+2*xb)*dx[1]*dydx[0] + dx[0]*dx[0]*dydx[1]) / xb
		for i := uint(1); i < (n - 1); i++ {
			d[i] = 3 * (dx[i]*dydx[i-1] + dx[i-1]*dydx[i])
		}
		d[n-1] = (dx[n-2]*dx[n-2]*dydx[n-3] + (2*xe+dx[n-2])*dx[n-3]*dydx[n-2]) / xe

		slopes := system.ComputeTridiagonal(a, b, c, d)

		s.breaks = make([]float64, n)
		copy(s.breaks, x)

		s.weights = make([]float64, 4*(n-1))
		for i := uint(0); i < (n - 1); i++ {
			α := (dydx[i] - slopes[i]) / dx[i]
			β := (slopes[i+1] - dydx[i]) / dx[i]
			s.weights[4*i+0] = (β - α) / dx[i]
			s.weights[4*i+1] = 2*α - β
			s.weights[4*i+2] = slopes[i]
			s.weights[4*i+3] = y[i]
		}
	}

	return s
}

// Compute calculates the ordinates corresponding to a series of abscissae. The
// abscissae are assumed to be an increasing sequence.
func (s *Cubic) Compute(x []float64) []float64 {
	n, nb := uint(len(x)), uint(len(s.breaks))-1

	y := make([]float64, n)

	for i, k := uint(0), uint(0); i < n; i++ {
		for x[i] > s.breaks[k+1] && k < (nb-1) {
			k++
		}

		z := x[i] - s.breaks[k]
		y[i] = s.weights[4*k]
		y[i] = z*y[i] + s.weights[4*k+1]
		y[i] = z*y[i] + s.weights[4*k+2]
		y[i] = z*y[i] + s.weights[4*k+3]
	}

	return y
}
