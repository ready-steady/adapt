package hybrid

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Continue gets called at the end of each iteration. If the function
	// returns false, the interpolation process is terminated. The first
	// argument is the set of currently active indices.
	Continue(*Active, *Progress) bool

	// Compute evaluates the target function at a point. The function is called
	// multiple times per iteration, depending on the number of active nodes.
	Compute(point, value []float64)

	// Score assigns a score to a location. The function is called after
	// Compute, and it is called as many times as Compute.
	Score(*Location) float64
}

// Location contains information about a dimensional location.
type Location struct {
	Values    []float64 // Target-function values
	Surpluses []float64 // Hierarchical surpluses
	Volumes   []float64 // Volumes under the basis functions
}

// Progress contains information about the interpolation process.
type Progress struct {
	More uint // Number of nodes to be evaluated
	Done uint // Number of nodes evaluated so far
}
