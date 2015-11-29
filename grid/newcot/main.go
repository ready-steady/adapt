// Package newcot provides means for working with the Newton–Cotes grid on the
// unit hypercube.
//
// Each node in the grid is identified by a sequence of levels and orders. Such
// a sequence is encoded as a sequence of uint64s where each uint64 is
// (level|order<<8). Consequently, the maximal level is 2^8, and the maximal
// order is 2^56.
package newcot

const (
	LEVEL_MASK = 0xFF
	LEVEL_SIZE = 8
	ORDER_SIZE = 64 - LEVEL_SIZE
)
