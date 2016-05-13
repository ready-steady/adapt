package internal

import (
	"math"

	"github.com/ready-steady/adapt/internal"
)

const (
	// None represents an absent choice.
	None = ^uint(0)
)

var (
	// Infinity is +∞.
	Infinity = math.Inf(1.0)
)

// Average returns the average value of a vector.
func Average(data []float64) (result float64) {
	for _, value := range data {
		result += value
	}
	result /= float64(len(data))
	return
}

// Choose selects an element with the highest positive priority.
func Choose(priority []float64, scope map[uint]bool) uint {
	if len(scope) == 0 {
		return None
	}
	k, max := None, 0.0
	for i := range scope {
		if k == None || priority[i] > max || (priority[i] == max && k > i) {
			k, max = i, priority[i]
		}
	}
	if max <= 0.0 {
		return None
	}
	return k
}

// Levelize returns the uniform norms of the levels of a set of indices.
func Levelize(indices []uint64, ni uint) (result []uint64) {
	nn := uint(len(indices)) / ni
	result = make([]uint64, nn)
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < ni; j++ {
			result[i] += internal.LEVEL_MASK & indices[i*ni+j]
		}
	}
	return
}

// MaxAbsolute returns the maximum absolute value of a vector.
func MaxAbsolute(data []float64) (result float64) {
	for i, n := uint(0), uint(len(data)); i < n; i++ {
		result = math.Max(result, math.Abs(data[i]))
	}
	return
}

// Set overwrites a vector with a fixed value.
func Set(data []float64, value float64) {
	for i := range data {
		data[i] = value
	}
}

// Subtract returns the difference between two vectors.
func Subtract(minuend, subtrahend []float64) (result []float64) {
	result = make([]float64, len(minuend))
	for i := range minuend {
		result[i] = minuend[i] - subtrahend[i]
	}
	return
}

// SumAbsolute returns the sum of the absolute values of a vector.
func SumAbsolute(data []float64) (result float64) {
	for _, value := range data {
		result += math.Abs(value)
	}
	return
}
