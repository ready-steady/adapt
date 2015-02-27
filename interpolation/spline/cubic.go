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
// should be a strictly increasing sequence with at least two elements. The
// ordinates can be multidimensional.
func NewCubic(x, y []float64) *Cubic {
	n := len(x)
	nd := len(y) / n

	dx := make([]float64, n-1)
	dydx := make([]float64, (n-1)*nd)
	for i := 0; i < (n - 1); i++ {
		dx[i] = x[i+1] - x[i]
		for j := 0; j < nd; j++ {
			dydx[i*nd+j] = (y[(i+1)*nd+j] - y[i*nd+j]) / dx[i]
		}
	}

	s := &Cubic{}

	switch n {
	case 2:
		s.breaks = []float64{x[0], x[1]}
		s.weights = make([]float64, nd*4)
		for j := 0; j < nd; j++ {
			s.weights[j*4+2] = dydx[j]
			s.weights[j*4+3] = y[j]
		}
	case 3:
		s.breaks = []float64{x[0], x[2]}
		s.weights = make([]float64, nd*4)
		for j := 0; j < nd; j++ {
			c1 := (dydx[nd+j] - dydx[j]) / (x[2] - x[0])
			s.weights[j*4+1] = c1
			s.weights[j*4+2] = dydx[j] - c1*dx[0]
			s.weights[j*4+3] = y[j]
		}
	default:
		xb := x[2] - x[0]
		xe := x[n-1] - x[n-3]

		a := make([]float64, n)
		for i := 0; i < (n - 2); i++ {
			a[i] = dx[i+1]
		}
		a[n-2] = xe

		b := make([]float64, n)
		b[0] = dx[1]
		for i := 1; i < (n - 1); i++ {
			b[i] = 2 * (dx[i] + dx[i-1])
		}
		b[n-1] = dx[n-3]

		c := make([]float64, n)
		c[1] = xb
		for i := 2; i < n; i++ {
			c[i] = dx[i-2]
		}

		d := make([]float64, nd*n)
		for j := 0; j < nd; j++ {
			d[j*n] = ((dx[0]+2*xb)*dx[1]*dydx[j] + dx[0]*dx[0]*dydx[nd+j]) / xb
			for i := 1; i < (n - 1); i++ {
				d[j*n+i] = 3 * (dx[i]*dydx[(i-1)*nd+j] + dx[i-1]*dydx[i*nd+j])
			}
			d[j*n+n-1] = (dx[n-2]*dx[n-2]*dydx[(n-3)*nd+j] +
				(2*xe+dx[n-2])*dx[n-3]*dydx[(n-2)*nd+j]) / xe
		}

		slopes := system.ComputeTridiagonal(a, b, c, d)

		s.breaks = make([]float64, n)
		copy(s.breaks, x)

		s.weights = make([]float64, (n-1)*nd*4)
		for i := 0; i < (n - 1); i++ {
			for j := 0; j < nd; j++ {
				α := (dydx[i*nd+j] - slopes[j*n+i]) / dx[i]
				β := (slopes[j*n+i+1] - dydx[i*nd+j]) / dx[i]
				s.weights[i*nd*4+j*4] = (β - α) / dx[i]
				s.weights[i*nd*4+j*4+1] = 2*α - β
				s.weights[i*nd*4+j*4+2] = slopes[j*n+i]
				s.weights[i*nd*4+j*4+3] = y[i*nd+j]
			}
		}
	}

	return s
}

// Compute calculates the ordinates corresponding to a series of abscissae. The
// abscissae should be an increasing sequence.
func (s *Cubic) Compute(x []float64) []float64 {
	n, nb := len(x), len(s.breaks)
	nd := len(s.weights) / (4 * (nb - 1))

	y := make([]float64, n*nd)

	for i, k := 0, 0; i < n; i++ {
		for x[i] > s.breaks[k+1] && k < (nb-2) {
			k++
		}

		z := x[i] - s.breaks[k]
		for j := 0; j < nd; j++ {
			y[i*nd+j] = s.weights[k*nd*4+j*4]
			y[i*nd+j] = z*y[i*nd+j] + s.weights[k*nd*4+j*4+1]
			y[i*nd+j] = z*y[i*nd+j] + s.weights[k*nd*4+j*4+2]
			y[i*nd+j] = z*y[i*nd+j] + s.weights[k*nd*4+j*4+3]
		}
	}

	return y
}
