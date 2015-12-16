// Package hybrid provides an algorithm for hierarchical interpolation with
// hybrid adaptation.
package hybrid

// Basis is a functional basis.
type Basis interface {
	// Compute evaluates the value of a basis function.
	Compute([]uint64, []float64) float64
}

// Grid is a sparse grid.
type Grid interface {
	// Compute returns the nodes corresponding to a set of indices.
	Compute([]uint64) []float64
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	basis Basis
	grid  Grid
}
