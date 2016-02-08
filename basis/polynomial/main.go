// Package polynomial provides functions for working with the basis formed by
// piecewise polynomial functions.
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
