package internal

import (
	"math"
	"sync"
)

var (
	infinity = math.Inf(1.0)
)

func Index(indexer Indexer, lindices []uint64, ni uint) ([]uint64, []uint) {
	nn := uint(len(lindices)) / ni
	indices, counts := []uint64(nil), make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		newIndices := indexer.Index(lindices[:ni])
		indices = append(indices, newIndices...)
		counts[i] = uint(len(newIndices)) / ni
		lindices = lindices[ni:]
	}
	return indices, counts
}

// Invoke evaluates a function at multiple nodes using multiple goroutines.
func Invoke(compute func([]float64, []float64), nodes []float64, ni, no, nw uint) []float64 {
	nn := uint(len(nodes)) / ni

	values := make([]float64, nn*no)

	jobs := make(chan uint, nn)
	group := sync.WaitGroup{}
	group.Add(int(nn))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				compute(nodes[j*ni:(j+1)*ni], values[j*no:(j+1)*no])
				group.Done()
			}
		}()
	}

	for i := uint(0); i < nn; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}

// LocateMaxFloat64s returns the position of the maximal element among a subset
// of a vector’s elements.
func LocateMaxFloat64s(data []float64, set map[uint]bool) uint {
	value, position := -infinity, ^uint(0)
	for i := range set {
		if data[i] > value {
			value, position = data[i], i
		}
	}
	return position
}

// LocateMinUint64s returns the position of the maximal element among a subset
// of a vector’s elements.
func LocateMinUint64s(data []uint64, set map[uint]bool) uint {
	value, position := ^uint64(0), ^uint(0)
	for i := range set {
		if data[i] < value {
			value, position = data[i], i
		}
	}
	return position
}

// MaxUint returns the maximal element among two.
func MaxUint(one uint, other uint) uint {
	if one > other {
		return one
	} else {
		return other
	}
}

// MaxUint64s returns the maximal element of a vector.
func MaxUint64s(data []uint64) uint64 {
	result := uint64(0)
	for _, value := range data {
		if value > result {
			result = value
		}
	}
	return result
}

// Subtract returns the difference between two vectors.
func Subtract(minuend, subtrahend []float64) []float64 {
	difference := make([]float64, len(minuend))
	for i := range minuend {
		difference[i] = minuend[i] - subtrahend[i]
	}
	return difference
}

// RepeatFloat64 returns a vector filled out with a particular value.
func RepeatFloat64(value float64, times uint) []float64 {
	data := make([]float64, times)
	for i := uint(0); i < times; i++ {
		data[i] = value
	}
	return data
}
