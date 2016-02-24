// Package external contains types and functions that are shared by the
// interpolation algorithm and appear in public interfaces.
package external

// Integrator computes the integral of a basis function.
type Integrator interface {
	Integrate([]uint64) float64
}

// Progress contains information about the interpolation process.
type Progress struct {
	More uint // Number of nodes to be evaluated
	Done uint // Number of nodes evaluated so far
}

// Set is a subset of ordered elements.
type Set map[uint]bool
