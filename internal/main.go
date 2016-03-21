// Package internal contains constants shared by the subpackages.
package internal

// An element of an nd-dimensional space is encoded by nd pairs (level, order).
// Each pair is a uint64 equal to (level|order<<LEVEL_SIZE) where LEVEL_SIZE is
// set to 6. In this encoding, the maximal level is 2^LEVEL_SIZE, and the
// maximal order is 2^(64-LEVEL_SIZE).
const (
	LEVEL_MASK = 0x3F
	LEVEL_SIZE = 6
	ORDER_SIZE = 64 - LEVEL_SIZE
)
