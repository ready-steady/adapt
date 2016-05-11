package algorithm

// Element contains information about an interpolation element.
type Element struct {
	Index []uint64 // Nodal index

	Node    []float64 // Grid node
	Volume  float64   // Basis-function volume
	Value   []float64 // Target-function value
	Surplus []float64 // Hierarchical surplus
}

// State contains information about an interpolation iteration.
type State struct {
	Ildices []uint64 // Level indices
	Indices []uint64 // Nodal indices
	Counts  []uint   // Number of nodal indices for each level index

	Nodes     []float64 // Grid nodes
	Volumes   []float64 // Basis-function volumes
	Values    []float64 // Target-function values
	Estimates []float64 // Approximated values
	Surpluses []float64 // Hierarchical surpluses
	Scores    []float64 // Nodal-index scores
}

// Strategy controls the interpolation process.
type Strategy interface {
	// First returns the initial state of the first iteration.
	First() *State

	// Next returns the initial state of the next iteration.
	Next(*State, *Surrogate) *State

	// Score assigns a score to an interpolation element.
	Score(*Element) float64
}
