package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/external"
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Before gets called once per iteration before involving Compute. If the
	// function returns false, the interpolation process is terminated.
	Before(*Progress) bool

	// Compute evaluates the target function at a point. The function is called
	// multiple times per iteration, depending on the number of active nodes.
	Compute(point, value []float64)

	// Score assigns a score to a location. The function is called after
	// Compute, and it is called as many times as Compute.
	Score(*Location) float64

	// After gets called once per iteration after involving Compute and Score.
	// The argument of the function is the set of currently active indices. If
	// the function returns false, the interpolation process is terminated.
	After(external.Set) bool
}

// Location contains information about a dimensional location.
type Location struct {
	Values    []float64 // Target-function values
	Surpluses []float64 // Hierarchical surpluses
}

// Progress contains information about the interpolation process.
type Progress struct {
	Level     uint // Reached level
	Active    uint // Number of active level indices
	Passive   uint // Number of passive level indices
	Requested uint // Number of requested function evaluations
	Performed uint // Number of performed function evaluations
}
