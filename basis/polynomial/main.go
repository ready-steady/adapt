// Package polynomial provides functions for working with the basis formed by
// piecewise polynomial functions.

// Each function of an nd-dimensional basis is given by nd pairs (level, order).
// Each pair is given as a uint64 equal to (level|order<<levelSize) where
// levelSize is set to 6. In this encoding, the maximal level is 2^levelSize,
// and the maximal order is 2^(64-levelSize).
package polynomial

import (
	"math"
)

const (
	levelMask = 0x3F
	levelSize = 6
	orderSize = 64 - levelSize
)

func equal(one, two float64) bool {
	const ε = 1e-14 // ~= 2^(-46)
	return one == two || math.Abs(one-two) < ε
}
