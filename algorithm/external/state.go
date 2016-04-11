package external

// State contains information about an interpolation iteration.
type State struct {
	Lindices []uint64 // Level indices
	Indices  []uint64 // Nodal indices
	Counts   []uint   // Number of nodal indices for each level index

	Nodes        []float64 // Grid nodes
	Volumes      []float64 // Basis-function volumes
	Observations []float64 // Target-function values
	Predictions  []float64 // Approximated values
	Surpluses    []float64 // Hierarchical surpluses
	Scores       []float64 // Nodal-index scores
}
