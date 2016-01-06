package hybrid

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

// String returns a human-friendly representation.
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
