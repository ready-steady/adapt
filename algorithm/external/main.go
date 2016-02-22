// Package external contains types and functions that are shared by the
// interpolation algorithm and appear in public interfaces.
package external

// Basis is a functional basis.
type Basis interface {
	// Compute evaluates the value of a basis function.
	Compute([]uint64, []float64) float64

	// Integrate computes the integral of a basis function.
	Integrate([]uint64) float64
}

// Set is a subset of ordered elements.
type Set map[uint]bool
