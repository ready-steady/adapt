package internal

import (
	"math"

	"github.com/ready-steady/adapt/internal"
)

var (
	infinity = math.Inf(1.0)
)

// Average returns the average value of a vector’s elements.
func Average(data []float64) float64 {
	return Sum(data) / float64(len(data))
}

// IsAdmissible checks if a set of indices is admissible.
func IsAdmissible(indices []uint64, ni uint) bool {
	nn := uint(len(indices)) / ni

	indices = append([]uint64(nil), indices...)
	for i := range indices {
		indices[i] = internal.LEVEL_MASK & indices[i]
	}

	hash := NewHash(ni)
	mapping := make(map[string]bool)
	for i := uint(0); i < nn; i++ {
		index := indices[i*ni : (i+1)*ni]
		mapping[hash.Key(index)] = true
	}

	for i := uint(0); i < nn; i++ {
		index := indices[i*ni : (i+1)*ni]
		for j := uint(0); j < ni; j++ {
			if index[j] == 0 {
				continue
			}
			index[j] -= 1
			_, ok := mapping[hash.Key(index)]
			index[j] += 1
			if !ok {
				return false
			}
		}
	}

	return true
}

// IsUnique checks if a set of indices has no repetitions.
func IsUnique(indices []uint64, ni uint) bool {
	unique := NewUnique(ni)

	indices = append([]uint64(nil), indices...)
	before := uint(len(indices)) / ni

	indices = unique.Distil(indices)
	after := uint(len(indices)) / ni

	return before == after
}

// LocateMax returns the position of the maximal element across a subset of a
// vector’s elements.
func LocateMax(data []float64, set map[uint]bool) uint {
	position, value := ^uint(0), 0.0
	for i := range set {
		if position == ^uint(0) || data[i] > value || (data[i] == value && position > i) {
			position, value = i, data[i]
		}
	}
	return position
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

// Max returns the maximum value across a vector’s elements.
func Max(data []float64) (result float64) {
	result = -infinity
	for i, n := uint(0), uint(len(data)); i < n; i++ {
		result = math.Max(result, data[i])
	}
	return
}

// MaxAbsolute returns the maximum absolute value across a vector’s elements.
func MaxAbsolute(data []float64) (result float64) {
	for i, n := uint(0), uint(len(data)); i < n; i++ {
		result = math.Max(result, math.Abs(data[i]))
	}
	return
}

// Subtract returns the difference between two vectors.
func Subtract(minuend, subtrahend []float64) (result []float64) {
	result = make([]float64, len(minuend))
	for i := range minuend {
		result[i] = minuend[i] - subtrahend[i]
	}
	return
}

// Sum returns the sum of a vector’s elements.
func Sum(data []float64) (result float64) {
	for _, value := range data {
		result += value
	}
	return
}

// SumAbsolute returns the sum of the absolute values of a vector’s elements.
func SumAbsolute(data []float64) (result float64) {
	for _, value := range data {
		result += math.Abs(value)
	}
	return
}
