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
	Inputs  uint
	Outputs uint

	Level uint
	Nodes uint

	Indices   []uint64
	Surpluses []float64
}

func (s *Surrogate) initialize(ni, no uint) {
	s.Inputs, s.Outputs, s.Nodes = ni, no, bufferInitCount

	s.Indices = make([]uint64, bufferInitCount*ni)
	s.Surpluses = make([]float64, bufferInitCount*no)
}

func (s *Surrogate) finalize(level uint, nn uint) {
	s.Level = level
	s.Nodes = nn

	s.Indices = s.Indices[0 : nn*s.Inputs]
	s.Surpluses = s.Surpluses[0 : nn*s.Outputs]
}

func (s *Surrogate) resize(nn uint) {
	if nn <= s.Nodes {
		return
	}

	if n := bufferGrowFactor * s.Nodes; n > nn {
		nn = n
	}

	indices := make([]uint64, nn*s.Inputs)
	surpluses := make([]float64, nn*s.Outputs)

	copy(indices, s.Indices[0:s.Nodes*s.Inputs])
	copy(surpluses, s.Surpluses[0:s.Nodes*s.Outputs])

	s.Nodes = nn

	s.Indices = indices
	s.Surpluses = surpluses
}

// String returns a string containing human-friendly information about the
// surrogate.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{inputs: %d, outputs: %d, level: %d, nodes: %d}",
		s.Inputs, s.Outputs, s.Level, s.Nodes)
}
