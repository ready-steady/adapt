// Package hybrid provides an algorithm for hierarchical interpolation with
// hybrid adaptation.
package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/global"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Basis is an interpolation basis.
type Basis interface {
	global.Basis
}

// Grid is an interpolation grid.
type Grid interface {
	global.Grid
	internal.GridRefinerToward
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	global.Interpolator
}

// New creates an interpolator.
func New(inputs, outputs uint, grid Grid, basis Basis) *Interpolator {
	return &Interpolator{*global.New(inputs, outputs, grid, basis)}
}
