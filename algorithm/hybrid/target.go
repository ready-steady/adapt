package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the target function at a point.
	Compute(point, value []float64)

	// Score assigns a score to a location.
	Score(*Location) float64

	// Done checks if the accuracy requirements have been satiated.
	Done(internal.Set) bool

	// Monitor gets called at the beginning of each iteration.
	Monitor(*Progress)
}

// Location contains information about a dimensional location.
type Location struct {
	Values    []float64 // Target-function values
	Surpluses []float64 // Hierarchical surpluses
}

// Progress contains information about the interpolation process.
type Progress struct {
	Level       uint // Reached level
	Active      uint // Number of active level indices
	Passive     uint // Number of passive level indices
	Evaluations uint // Number of function evaluations
}
