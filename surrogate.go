package adapt

import (
	"fmt"
)

// Surrogate is an interpolant for a function.
type Surrogate struct {
	Inputs    uint      // Number of inputs
	Outputs   uint      // Number of outputs
	Level     uint      // Level of interpolation
	Nodes     uint      // Number of nodes
	Indices   []uint64  // Indices of the nodes
	Surpluses []float64 // Hierarchical surpluses
	Accept    []uint    // Number of nodes accepted at each iteration
	Reject    []uint    // Number of nodes rejected at each iteration
}

func newSurrogate(ni, no uint) *Surrogate {
	return &Surrogate{
		Inputs:    ni,
		Outputs:   no,
		Indices:   make([]uint64, 0, ni),
		Surpluses: make([]float64, 0, no),
		Accept:    make([]uint, 0),
		Reject:    make([]uint, 0),
	}
}

func (s *Surrogate) push(indices []uint64, surpluses []float64) {
	s.Indices = append(s.Indices, indices...)
	s.Surpluses = append(s.Surpluses, surpluses...)
}

func (s *Surrogate) step(level, accepted, rejected uint) {
	s.Level = level
	s.Nodes += accepted
	s.Accept = append(s.Accept, accepted)
	s.Reject = append(s.Reject, rejected)
}

// String returns human-friendly information about the interpolant.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{inputs: %d, outputs: %d, level: %d, nodes: %d}",
		s.Inputs, s.Outputs, s.Level, s.Nodes)
}
