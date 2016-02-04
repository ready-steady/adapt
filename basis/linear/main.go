// Package linear provides functions for working with the basis formed by
// piecewise linear functions.
//
// Each function in the basis is identified by a sequence of levels and orders.
// Such a sequence is encoded as a sequence of uint64s where each uint64 is
// (level|order<<6). Consequently, the maximal level is 2^6, and the maximal
// order is 2^58.
package linear

const (
	levelMask = 0x3F
	levelSize = 6
	orderSize = 64 - levelSize
)
