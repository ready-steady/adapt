// Package external contains types and functions that are shared by the
// interpolation algorithm and appear in public interfaces.
package external

// Integrator computes the integral of a basis function.
type Integrator interface {
	Integrate([]uint64) float64
}
