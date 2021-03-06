package algorithm

import (
	"fmt"
)

// Surrogate is an interpolant for a function.
type Surrogate struct {
	Inputs  uint // Number of inputs
	Outputs uint // Number of outputs
	Nodes   uint // Number of nodes

	Indices   []uint64  // Indices of the nodes
	Surpluses []float64 // Hierarchical surpluses
	Integral  []float64 // Integral over the whole domain
}

// NewSurrogate returns an empty surrogate.
func NewSurrogate(ni, no uint) *Surrogate {
	return &Surrogate{
		Inputs:  ni,
		Outputs: no,

		Indices:   make([]uint64, 0),
		Surpluses: make([]float64, 0),
		Integral:  make([]float64, no),
	}
}

// Push takes into account new indices and surpluses.
func (self *Surrogate) Push(indices []uint64, surpluses, volumes []float64) {
	self.Nodes += uint(len(indices)) / self.Inputs
	self.Indices = append(self.Indices, indices...)
	self.Surpluses = append(self.Surpluses, surpluses...)
	cumulate(indices, surpluses, volumes, self.Inputs, self.Outputs, self.Integral)
}

// String returns a summary.
func (self *Surrogate) String() string {
	phantom := struct {
		inputs  uint
		outputs uint
		nodes   uint
	}{
		inputs:  self.Inputs,
		outputs: self.Outputs,
		nodes:   self.Nodes,
	}
	return fmt.Sprintf("%+v", phantom)
}

func cumulate(indices []uint64, surpluses, volumes []float64, ni, no uint, integral []float64) {
	nn := uint(len(indices)) / ni
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < no; j++ {
			integral[j] += surpluses[i*no+j] * volumes[i]
		}
	}
}
