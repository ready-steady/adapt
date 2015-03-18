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

func newSurrogate(ni, no uint) *Surrogate {
	return &Surrogate{
		Inputs:  ni,
		Outputs: no,

		Nodes: bufferInitCount,

		Indices:   make([]uint64, bufferInitCount*ni),
		Surpluses: make([]float64, bufferInitCount*no),
	}
}

func (s *Surrogate) finalize(nn uint) {
	ni, no := s.Inputs, s.Outputs

	s.Indices = s.Indices[:nn*ni]
	s.Surpluses = s.Surpluses[:nn*no]

	level := uint(0)
	for i := uint(0); i < nn; i++ {
		l := uint(0)
		for j := uint(0); j < ni; j++ {
			l += uint(0xFFFFFFFF & s.Indices[i*ni+j])
		}
		if l > level {
			level = l
		}
	}

	s.Level = level
	s.Nodes = nn
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

	copy(indices, s.Indices[:s.Nodes*s.Inputs])
	copy(surpluses, s.Surpluses[:s.Nodes*s.Outputs])

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
