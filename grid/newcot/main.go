// Package newcot provides means for working with the Newton–Cotes grid on the
// unit hypercube.
//
// Each node in the grid is identified by a sequence of levels and orders.
// Throughout the package, such a sequence is encoded as a sequence of uint64s,
// referred to as an index, where each uint64 is (level|order<<32).
package newcot
