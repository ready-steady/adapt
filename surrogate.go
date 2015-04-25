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

	Steps     []uint
	Indices   []uint64
	Surpluses []float64
}

func newSurrogate(ni, no uint) *Surrogate {
	return &Surrogate{
		Inputs:  ni,
		Outputs: no,

		Steps:     make([]uint, 0),
		Indices:   make([]uint64, 0, ni),
		Surpluses: make([]float64, 0, no),
	}
}

func (s *Surrogate) push(indices []uint64, surpluses []float64) {
	s.Steps = append(s.Steps, uint(len(indices))/s.Inputs)
	s.Indices = append(s.Indices, indices...)
	s.Surpluses = append(s.Surpluses, surpluses...)
}

// String returns a string containing human-friendly information about the
// surrogate.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{inputs: %d, outputs: %d, level: %d, nodes: %d}",
		s.Inputs, s.Outputs, s.Level, s.Nodes)
}
