package internal

import (
	"sync"
)

// Approximate evaluates an interpolant at multiple points using multiple
// goroutines.
func Approximate(basis Basis, indices []uint64, surpluses, points []float64,
	ni, no, nw uint) []float64 {

	nn := uint(len(indices)) / ni
	np := uint(len(points)) / ni
	values := make([]float64, np*no)

	jobs := make(chan uint, np)
	group := sync.WaitGroup{}
	group.Add(int(np))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				point := points[j*ni : (j+1)*ni]
				value := values[j*no : (j+1)*no]

				for k := uint(0); k < nn; k++ {
					weight := basis.Compute(indices[k*ni:(k+1)*ni], point)
					if weight == 0.0 {
						continue
					}
					for l := uint(0); l < no; l++ {
						value[l] += weight * surpluses[k*no+l]
					}
				}

				group.Done()
			}
		}()
	}

	for i := uint(0); i < np; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
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

// Subtract returns the difference between two vectors.
func Subtract(minuend, subtrahend []float64) []float64 {
	difference := make([]float64, len(minuend))
	for i := range minuend {
		difference[i] = minuend[i] - subtrahend[i]
	}
	return difference
}

// MaxUint64 returns the maximal element.
func MaxUint64(data []uint64) uint64 {
	result := uint64(0)
	for _, value := range data {
		if value > result {
			result = value
		}
	}
	return result
}

func maxFloat64Set(data []float64, set Set) (float64, uint) {
	value, position := -infinity, ^uint(0)
	for i := range set {
		if data[i] > value {
			value, position = data[i], i
		}
	}
	return value, position
}

func minUint64Set(data []uint64, set Set) (uint64, uint) {
	value, position := ^uint64(0), ^uint(0)
	for i := range set {
		if data[i] < value {
			value, position = data[i], i
		}
	}
	return value, position
}
