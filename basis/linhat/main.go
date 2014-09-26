// Package linhat provides functions for working with the basis formed by the
// linear hat function on the unit hypercube.
//
// Each function in the basis is identified by a sequence of levels and orders.
// Throughout the package, such a sequence is encoded as a sequence of uint64s,
// referred to as an index, where each uint64 is (level|order<<32).
package linhat
