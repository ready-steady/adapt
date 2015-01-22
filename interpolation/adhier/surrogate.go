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

	ic uint32
	oc uint32
	nc uint32

	indices   []uint64
	surpluses []float64
}

func (s *Surrogate) initialize(ic, oc uint32) {
	s.ic, s.oc, s.nc = ic, oc, bufferInitCount

	s.indices = make([]uint64, bufferInitCount*ic)
	s.surpluses = make([]float64, bufferInitCount*oc)
}

func (s *Surrogate) finalize(level uint8, nc uint32) {
	s.level = level
	s.nc = nc

	s.indices = s.indices[0 : nc*s.ic]
	s.surpluses = s.surpluses[0 : nc*s.oc]
}

func (s *Surrogate) resize(nc uint32) {
	if nc <= s.nc {
		return
	}

	if count := bufferGrowFactor * s.nc; count > nc {
		nc = count
	}

	indices := make([]uint64, nc*s.ic)
	surpluses := make([]float64, nc*s.oc)

	copy(indices, s.indices[0:s.nc*s.ic])
	copy(surpluses, s.surpluses[0:s.nc*s.oc])

	s.nc = nc

	s.indices = indices
	s.surpluses = surpluses
}

// String returns a string containing human-friendly information about the
// surrogate.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{inputs: %d, outputs: %d, level: %d, nodes: %d}",
		s.ic, s.oc, s.level, s.nc)
}
