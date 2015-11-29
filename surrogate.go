package adapt

import (
	"fmt"
)

// Surrogate is an interpolant for a function.
type Surrogate struct {
	Inputs    uint      // Number of inputs
	Outputs   uint      // Number of outputs
	Level     uint      // Interpolation level
	Nodes     uint      // Number of nodes
	Indices   []uint64  // Indices of the nodes
	Surpluses []float64 // Hierarchical surpluses of the nodes
	Passive   []uint    // Number of passive nodes at each iteration
}

func newSurrogate(ni, no uint) *Surrogate {
	return &Surrogate{
		Inputs:    ni,
		Outputs:   no,
		Indices:   make([]uint64, 0, ni),
		Surpluses: make([]float64, 0, no),
		Passive:   make([]uint, 0),
	}
}

func (self *Surrogate) push(indices []uint64, surpluses []float64) {
	self.Indices = append(self.Indices, indices...)
	self.Surpluses = append(self.Surpluses, surpluses...)
}

func (self *Surrogate) step(level, passive uint) {
	self.Level = level
	self.Nodes += passive
	self.Passive = append(self.Passive, passive)
}

// String returns human-friendly information about the interpolant.
func (self *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{inputs: %d, outputs: %d, level: %d, nodes: %d}",
		self.Inputs, self.Outputs, self.Level, self.Nodes)
}
