// Package equidistant provides means for working with the Newton–Cotes grid.
//
// Each node of an nd-dimensional grid is given by nd pairs (level, order). Each
// pair is given as a uint64 equal to (level|order<<levelSize) where levelSize
// is set to 6. In this encoding, the maximal level is 2^levelSize, and the
// maximal order is 2^(64-levelSize).
package equidistant

const (
	levelMask = 0x3F
	levelSize = 6
	orderSize = 64 - levelSize
)
