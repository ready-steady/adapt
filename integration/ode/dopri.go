package ode

import (
	"errors"
	"math"
)

// DormandPrince is an integrator based on the Dormand–Prince method.
//
// https://en.wikipedia.org/wiki/Dormand–Prince_method
type DormandPrince struct {
	config Config
}

// NewDormandPrince creates a new Dormand–Prince integrator.
func NewDormandPrince(config *Config) (*DormandPrince, error) {
	if err := config.verify(); err != nil {
		return nil, err
	}
	return &DormandPrince{config: *config}, nil
}

// Compute integrates the system of differential equations y' = f(x, y) and
// returns the resulting solution at the specified points. The derivative
// function is supposed to evaluate f(x, y) given x and y in its first and
// second arguments, respectively, and to store the computed value in its third
// argument. The initial value of y is given by initial, which corresponds to
// the first point in points.
func (self *DormandPrince) Compute(derivative func(float64, []float64, []float64),
	points []float64, initial []float64) ([]float64, error) {

	const (
		c2 = 1.0 / 5
		c3 = 3.0 / 10
		c4 = 4.0 / 5
		c5 = 8.0 / 9

		a21 = 1.0 / 5
		a31 = 3.0 / 40
		a32 = 9.0 / 40
		a41 = 44.0 / 45
		a42 = -56.0 / 15
		a43 = 32.0 / 9
		a51 = 19372.0 / 6561
		a52 = -25360.0 / 2187
		a53 = 64448.0 / 6561
		a54 = -212.0 / 729
		a61 = 9017.0 / 3168
		a62 = -355.0 / 33
		a63 = 46732.0 / 5247
		a64 = 49.0 / 176
		a65 = -5103.0 / 18656
		a71 = 35.0 / 384
		a73 = 500.0 / 1113
		a74 = 125.0 / 192
		a75 = -2187.0 / 6784
		a76 = 11.0 / 84

		e1 = 71.0 / 57600
		e3 = -71.0 / 16695
		e4 = 71.0 / 1920
		e5 = -17253.0 / 339200
		e6 = 22.0 / 525
		e7 = -1.0 / 40

		power = 1.0 / 5
	)

	pc := uint(len(points))
	if pc < 2 {
		return nil, errors.New("need at least two points")
	}

	dc := uint(len(initial))
	if dc == 0 {
		return nil, errors.New("need an initial value")
	}

	z := make([]float64, dc)

	y := make([]float64, dc)

	ynew := make([]float64, dc)

	f := make([]float64, 7*dc)

	f1 := f[0*dc : 1*dc]
	f2 := f[1*dc : 2*dc]
	f3 := f[2*dc : 3*dc]
	f4 := f[3*dc : 4*dc]
	f5 := f[4*dc : 5*dc]
	f6 := f[5*dc : 6*dc]
	f7 := f[6*dc : 7*dc]

	x, xend := points[0], points[pc-1]
	copy(y, initial)
	derivative(x, y, f1)

	values := make([]float64, pc*dc)
	copy(values, initial)
	cc := uint(1)

	config := &self.config

	abstol, reltol := config.AbsoluteTolerance, config.RelativeTolerance
	threshold := abstol / reltol

	// Compute the limits on the step size.
	hmin := 16 * epsilon(0)
	hmax := config.MaximalStep
	if hmax == 0 {
		hmax = 0.1 * (xend - x)
	}

	// Choose the initial step size.
	h := config.InitialStep
	if h == 0 {
		h = math.Min(hmax, points[1]-x)

		scale := math.Inf(-1)
		for i := uint(0); i < dc; i++ {
			scale = math.Max(scale, math.Abs(f1[i]/math.Max(math.Abs(y[i]), threshold)))
		}
		scale = scale / (0.8 * math.Pow(reltol, power))

		if h*scale > 1 {
			h = 1 / scale
		}

		h = math.Max(hmin, h)
	} else {
		h = math.Min(hmax, math.Max(hmin, h))
	}

	var xnew, ε float64
	var rejected bool

	for done := false; !done; {
		hmin := 16 * epsilon(x)
		h = math.Min(hmax, math.Max(hmin, h))

		// Close to the end?
		if 1.1*h >= xend-x {
			h = xend - x
			done = true
		}

		for rejected = false; ; {
			// Step 1
			for j := uint(0); j < dc; j++ {
				z[j] = y[j] + h*a21*f1[j]
			}

			// Step 2
			derivative(x+c2*h, z, f2)
			for j := uint(0); j < dc; j++ {
				z[j] = y[j] + h*(a31*f1[j]+a32*f2[j])
			}

			// Step 3
			derivative(x+c3*h, z, f3)
			for j := uint(0); j < dc; j++ {
				z[j] = y[j] + h*(a41*f1[j]+a42*f2[j]+a43*f3[j])
			}

			// Step 4
			derivative(x+c4*h, z, f4)
			for j := uint(0); j < dc; j++ {
				z[j] = y[j] + h*(a51*f1[j]+a52*f2[j]+a53*f3[j]+a54*f4[j])
			}

			// Step 5
			derivative(x+c5*h, z, f5)
			for j := uint(0); j < dc; j++ {
				z[j] = y[j] + h*(a61*f1[j]+a62*f2[j]+a63*f3[j]+a64*f4[j]+a65*f5[j])
			}

			// Step 6
			derivative(x+h, z, f6)
			for j := uint(0); j < dc; j++ {
				ynew[j] = y[j] + h*(a71*f1[j]+a73*f3[j]+a74*f4[j]+a75*f5[j]+a76*f6[j])
			}

			xnew = x + h
			h = xnew - x

			// Step 1
			derivative(xnew, ynew, f7)

			// Error control
			ε = 0
			for j := uint(0); j < dc; j++ {
				scale := math.Max(math.Max(math.Abs(y[j]), math.Abs(ynew[j])), threshold)
				Δ := e1*f1[j] + e3*f3[j] + e4*f4[j] + e5*f5[j] + e6*f6[j] + e7*f7[j]
				ε = math.Max(h*math.Abs(Δ/scale), ε)
			}

			if ε <= reltol {
				break
			}

			if h <= hmin {
				return nil, errors.New("encountered a step-size underflow")
			}

			if rejected {
				h = math.Max(hmin, 0.5*h)
			} else {
				h = math.Max(hmin, h*math.Max(0.1, 0.8*math.Pow(reltol/ε, power)))
				rejected = true
			}

			done = false
		}

		for cc < pc {
			if xnew-points[cc] < 0 {
				break
			}

			if points[cc] == xnew {
				copy(values[cc*dc:(cc+1)*dc], ynew)
			} else {
				self.interpolate(x, y, f, h, points[cc], values[cc*dc:(cc+1)*dc])
			}

			cc++
		}

		if done {
			break
		}

		if !rejected {
			if scale := 1.25 * math.Pow(ε/reltol, power); scale > 0.2 {
				h = h / scale
			} else {
				h = 5 * h
			}
		}

		x = xnew
		copy(f1, f7)
		copy(y, ynew)
	}

	return values, nil
}

func (_ *DormandPrince) interpolate(x float64, y, f []float64, h, xnext float64, ynext []float64) {
	const (
		c11 = 1.0
		c12 = -183.0 / 64
		c13 = 37.0 / 12
		c14 = -145.0 / 128
		c32 = 1500.0 / 371
		c33 = -1000.0 / 159
		c34 = 1000.0 / 371
		c42 = -125.0 / 32
		c43 = 125.0 / 12
		c44 = -375.0 / 64
		c52 = 9477.0 / 3392
		c53 = -729.0 / 106
		c54 = 25515.0 / 6784
		c62 = -11.0 / 7
		c63 = 11.0 / 3
		c64 = -55.0 / 28
		c72 = 3.0 / 2
		c73 = -4.0
		c74 = 5.0 / 2
	)

	dc := uint(len(y))

	s1 := (xnext - x) / h
	s2 := s1 * s1
	s3 := s1 * s2
	s4 := s1 * s3

	for i := uint(0); i < dc; i++ {
		f1 := f[0*dc+i]
		f3 := f[2*dc+i]
		f4 := f[3*dc+i]
		f5 := f[4*dc+i]
		f6 := f[5*dc+i]
		f7 := f[6*dc+i]

		ynext[i] = y[i] +
			h*s1*(c11*f1) +
			h*s2*(c12*f1+c32*f3+c42*f4+c52*f5+c62*f6+c72*f7) +
			h*s3*(c13*f1+c33*f3+c43*f4+c53*f5+c63*f6+c73*f7) +
			h*s4*(c14*f1+c34*f3+c44*f4+c54*f5+c64*f6+c74*f7)
	}
}

func epsilon(x float64) float64 {
	x = math.Abs(x)
	return math.Nextafter(x, x+1) - x
}
