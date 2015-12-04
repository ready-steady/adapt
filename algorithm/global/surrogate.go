package global

import (
	"fmt"
)

// Surrogate is an interpolant for a function.
type Surrogate struct {
	Inputs  uint // Number of inputs
	Outputs uint // Number of outputs
	Nodes   uint // Number of nodes
}

func newSurrogate(ni, no uint) *Surrogate {
	return &Surrogate{
		Inputs:  ni,
		Outputs: no,
	}
}

// String returns human-friendly information about the surrogate.
func (self *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{inputs: %d, outputs: %d, nodes: %d}",
		self.Inputs, self.Outputs, self.Nodes)
}
