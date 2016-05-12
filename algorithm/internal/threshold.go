package internal

import (
	"math"
)

// Threshold is an adaptive error threshold.
type Threshold struct {
	Values []float64 // Threshold

	lower []float64
	upper []float64

	no uint
	εa float64
	εr float64
}

// NewThreshold creates a Threshold.
func NewThreshold(outputs uint, absolute, relative float64) *Threshold {
	return &Threshold{
		Values: make([]float64, outputs),

		lower: repeat(math.Inf(1.0), outputs),
		upper: repeat(math.Inf(-1.0), outputs),

		no: outputs,
		εa: absolute,
		εr: relative,
	}
}

// Compress compresses multiple points into a single one so that it can later on
// be tested against the threshold.
func (self *Threshold) Compress(destination, source []float64) {
	no := self.no
	nn := uint(len(source)) / no
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < no; j++ {
			destination[j] = math.Max(destination[j], math.Abs(source[i*no+j]))
		}
	}
}

// Update updates the threshold.
func (self *Threshold) Update(values []float64) {
	no := self.no
	for i, m := uint(0), uint(len(values)); i < m; i++ {
		j := i % no
		self.lower[j] = math.Min(self.lower[j], values[i])
		self.upper[j] = math.Max(self.upper[j], values[i])
	}
	for i := uint(0); i < no; i++ {
		self.Values[i] = math.Max(self.εa, self.εr*(self.upper[i]-self.lower[i]))
	}
}

func repeat(value float64, times uint) []float64 {
	values := make([]float64, times)
	for i := uint(0); i < times; i++ {
		values[i] = value
	}
	return values
}
