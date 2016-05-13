package internal

import (
	"math"
)

// Threshold is an adaptive error threshold.
type Threshold struct {
	values []float64
	lower  []float64
	upper  []float64

	no uint
	εa float64
	εr float64
}

// NewThreshold creates a Threshold.
func NewThreshold(outputs uint, absolute, relative float64) *Threshold {
	return &Threshold{
		values: make([]float64, outputs),
		lower:  repeat(Infinity, outputs),
		upper:  repeat(-Infinity, outputs),

		no: outputs,
		εa: absolute,
		εr: relative,
	}
}

// Check checks if the threshold is satisfied.
func (self *Threshold) Check(errors []float64, include map[uint]bool) bool {
	no := self.no
	for i := range include {
		for j := uint(0); j < no; j++ {
			if errors[i*no+j] > self.values[j] {
				return false
			}
		}
	}
	return true
}

// Compress compresses multiple errors into a single one so that it can later on
// be tested against the threshold.
func (self *Threshold) Compress(error, errors []float64) {
	for i, m, no := uint(0), uint(len(errors)), self.no; i < m; i++ {
		j := i % no
		error[j] = math.Max(error[j], math.Abs(errors[i]))
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
		self.values[i] = math.Max(self.εa, self.εr*(self.upper[i]-self.lower[i]))
	}
}

func repeat(value float64, times uint) []float64 {
	values := make([]float64, times)
	for i := uint(0); i < times; i++ {
		values[i] = value
	}
	return values
}
