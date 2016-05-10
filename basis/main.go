// Package basis contains code shared by the interpolation bases.
package basis

// Computer returns the value of a basis function.
type Computer interface {
	Compute([]uint64, []float64) float64
}

// Integrator returns the integral of a basis function.
type Integrator interface {
	Integrate([]uint64) float64
}
