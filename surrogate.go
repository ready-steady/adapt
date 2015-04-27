package adapt

import (
	"fmt"
)

// Surrogate is the result of Compute, which represents an interpolant for a
// function.
type Surrogate struct {
	Inputs  uint
	Outputs uint
	Level   uint
	Nodes   uint

	Accept    []uint
	Reject    []uint
	Indices   []uint64
	Surpluses []float64
}

func newSurrogate(ni, no uint) *Surrogate {
	return &Surrogate{
		Inputs:  ni,
		Outputs: no,

		Accept:    make([]uint, 0),
		Reject:    make([]uint, 0),
		Indices:   make([]uint64, 0, ni),
		Surpluses: make([]float64, 0, no),
	}
}

func (s *Surrogate) push(indices []uint64, surpluses []float64, discard []bool) {
	ni, no := s.Inputs, s.Outputs
	nn, na := uint(len(discard)), uint(0)

	for i := uint(0); i < nn; i++ {
		if discard[i] {
			continue
		}
		na++
		s.Indices = append(s.Indices, indices[i*ni:(i+1)*ni]...)
		s.Surpluses = append(s.Surpluses, surpluses[i*no:(i+1)*no]...)
	}

	s.Accept = append(s.Accept, na)
	s.Reject = append(s.Reject, nn-na)
}

// String returns a string containing human-friendly information about the
// surrogate.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{inputs: %d, outputs: %d, level: %d, nodes: %d}",
		s.Inputs, s.Outputs, s.Level, s.Nodes)
}
