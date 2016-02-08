// Package polynomial provides functions for working with the basis formed by
// piecewise polynomial functions.
package polynomial

const (
	levelMask = 0x3F
	levelSize = 6
	orderSize = 64 - levelSize
)
