// Package internal contains types and functions that are shared by the
// interpolation algorithm and are used only internally.
package internal

// Basis is a functional basis.
type Basis interface {
	// Compute evaluates the value of a basis function.
	Compute([]uint64, []float64) float64
}
