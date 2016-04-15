// Package hybrid provides an algorithm for hierarchical interpolation with
// hybrid adaptation.
package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/global"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Algorithm is the interpolation algorithm.
type Algorithm struct {
	global.Algorithm
}

// Basis is an interpolation basis.
type Basis interface {
	global.Basis
}

// Grid is an interpolation grid.
type Grid interface {
	global.Grid
	internal.GridRefinerToward
}

// New creates an interpolator.
func New(inputs, outputs uint, grid Grid, basis Basis) *Algorithm {
	return &Algorithm{*global.New(inputs, outputs, grid, basis)}
}
