package internal

import (
	"math"
)

// Domain is a structure for keeping track of the domain of a function.
type Domain struct {
	Lower  []float64 // Lower bound
	Upper  []float64 // Upper bound
	Spread []float64 // Distance between the bounds

	no uint
}

// NewDomain creates a Domain.
func NewDomain(outputs uint) *Domain {
	return &Domain{
		Lower:  repeat(math.Inf(1.0), outputs),
		Upper:  repeat(math.Inf(-1.0), outputs),
		Spread: repeat(math.Inf(1.0), outputs),

		no: outputs,
	}
}

func (self *Domain) Update(values []float64) {
	no := self.no
	for i, m := uint(0), uint(len(values)); i < m; i++ {
		j := i % no
		self.Lower[j] = math.Min(self.Lower[j], values[i])
		self.Upper[j] = math.Max(self.Upper[j], values[i])
	}
	for i := uint(0); i < no; i++ {
		self.Spread[i] = self.Upper[i] - self.Lower[i]
	}
}

func repeat(value float64, times uint) []float64 {
	values := make([]float64, times)
	for i := uint(0); i < times; i++ {
		values[i] = value
	}
	return values
}
