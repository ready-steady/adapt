// Package linhat provides functions for working with the basis formed by the
// linear hat function on the unit hypercube.
//
// Each function in the basis is identified by a sequence of levels and orders.
// Such a sequence is encoded as a sequence of uint64s where each uint64 is
// (level|order<<6). Consequently, the maximal level is 2^6, and the maximal
// order is 2^58.
package linhat

const (
	LEVEL_MASK = 0x3F
	LEVEL_SIZE = 6
	ORDER_SIZE = 64 - LEVEL_SIZE
)
