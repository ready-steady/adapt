package external

// State contains information about an interpolation iteration.
type State struct {
	Indices []uint64 // Nodal indices

	Nodes        []float64 // Grid nodes
	Volumes      []float64 // Basis-function volumes
	Observations []float64 // Target-function values
	Predictions  []float64 // Approximated values
	Surpluses    []float64 // Hierarchical surpluses
	Scores       []float64 // Nodal-index scores
}

// FineState contains information about an interpolation iteration.
type FineState struct {
	Lindices []uint64 // Level indices
	Counts   []uint   // Number of nodal indices for each level index

	State
}
