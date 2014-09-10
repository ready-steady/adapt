package adhier

import (
	"fmt"
)

const (
	bufferInitCount  = 200
	bufferGrowFactor = 2
)

// Surrogate is the result of Compute, which represents an interpolant for a
// function.
type Surrogate struct {
	level uint8

	ic uint16
	oc uint16
	nc uint32

	levels    []uint8
	orders    []uint32
	surpluses []float64
}

func (s *Surrogate) initialize(ic, oc uint16) {
	is := bufferInitCount * uint32(ic)
	os := bufferInitCount * uint32(oc)

	s.ic = ic
	s.oc = oc
	s.nc = bufferInitCount

	s.levels = make([]uint8, is)
	s.orders = make([]uint32, is)
	s.surpluses = make([]float64, os)
}

func (s *Surrogate) finalize(level uint8, nc uint32) {
	is := nc * uint32(s.ic)
	os := nc * uint32(s.oc)

	s.level = level
	s.nc = nc

	s.levels = s.levels[0:is]
	s.orders = s.orders[0:is]
	s.surpluses = s.surpluses[0:os]
}

func (s *Surrogate) resize(nc uint32) {
	if nc <= s.nc {
		return
	}

	if count := bufferGrowFactor * s.nc; count > nc {
		nc = count
	}

	// New sizes
	is := nc * uint32(s.ic)
	os := nc * uint32(s.oc)

	levels := make([]uint8, is)
	orders := make([]uint32, is)
	surpluses := make([]float64, os)

	// Old sizes
	is = s.nc * uint32(s.ic)
	os = s.nc * uint32(s.oc)

	copy(levels, s.levels[0:is])
	copy(orders, s.orders[0:is])
	copy(surpluses, s.surpluses[0:os])

	s.nc = nc

	s.levels = levels
	s.orders = orders
	s.surpluses = surpluses
}

// String returns a string containing human-friendly information about the
// surrogate.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{ inputs: %d, outputs: %d, levels: %d, nodes: %d }",
		s.ic, s.oc, s.level, s.nc)
}
