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
	Inputs  uint32
	Outputs uint32

	Level uint8
	Nodes uint32

	Indices   []uint64
	Surpluses []float64
}

func (s *Surrogate) initialize(ic, oc uint32) {
	s.Inputs, s.Outputs, s.Nodes = ic, oc, bufferInitCount

	s.Indices = make([]uint64, bufferInitCount*ic)
	s.Surpluses = make([]float64, bufferInitCount*oc)
}

func (s *Surrogate) finalize(level uint8, nc uint32) {
	s.Level = level
	s.Nodes = nc

	s.Indices = s.Indices[0 : nc*s.Inputs]
	s.Surpluses = s.Surpluses[0 : nc*s.Outputs]
}

func (s *Surrogate) resize(nc uint32) {
	if nc <= s.Nodes {
		return
	}

	if count := bufferGrowFactor * s.Nodes; count > nc {
		nc = count
	}

	indices := make([]uint64, nc*s.Inputs)
	surpluses := make([]float64, nc*s.Outputs)

	copy(indices, s.Indices[0:s.Nodes*s.Inputs])
	copy(surpluses, s.Surpluses[0:s.Nodes*s.Outputs])

	s.Nodes = nc

	s.Indices = indices
	s.Surpluses = surpluses
}

// String returns a string containing human-friendly information about the
// surrogate.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{inputs: %d, outputs: %d, level: %d, nodes: %d}",
		s.Inputs, s.Outputs, s.Level, s.Nodes)
}
