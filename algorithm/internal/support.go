package internal

import (
	"math"
	"sync"

	"github.com/ready-steady/adapt/internal"
)

var (
	infinity = math.Inf(1.0)
)

// Index returns the nodal indices of a set of level indices.
func Index(indexer Indexer, lindices []uint64, ni uint) ([]uint64, []uint) {
	nn := uint(len(lindices)) / ni
	indices, counts := []uint64(nil), make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		newIndices := indexer.Index(lindices[i*ni : (i+1)*ni])
		indices = append(indices, newIndices...)
		counts[i] = uint(len(newIndices)) / ni
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

	indices = append([]uint64{}, indices...)
	before := uint(len(indices)) / ni

	indices = unique.Distil(indices)
	after := uint(len(indices)) / ni

	return before == after
}

// LocateMax returns the position of the maximal element among a subset of a
// vector’s elements.
func LocateMax(data []float64, set map[uint]bool) uint {
	value, position := -infinity, ^uint(0)
	for i := range set {
		if data[i] > value {
			value, position = data[i], i
		}
	}
	return position
}

// Levelize returns the uniform norms of the levels of a set of indices.
func Levelize(indices []uint64, ni uint) []uint64 {
	nn := uint(len(indices)) / ni
	levels := make([]uint64, nn)
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < ni; j++ {
			levels[i] += internal.LEVEL_MASK & indices[i*ni+j]
		}
	}
	return levels
}

// Subtract returns the difference between two vectors.
func Subtract(minuend, subtrahend []float64) []float64 {
	difference := make([]float64, len(minuend))
	for i := range minuend {
		difference[i] = minuend[i] - subtrahend[i]
	}
	return difference
}

// Repeat returns a vector filled out with a particular value.
func Repeat(value float64, times uint) []float64 {
	data := make([]float64, times)
	for i := uint(0); i < times; i++ {
		data[i] = value
	}
	return data
}
