// Package internal contains types and functions that are shared by the
// interpolation algorithm and are used only internally.
package internal

import (
	"sync"
)

// Computer computes the value of a basis function.
type Computer interface {
	Compute([]uint64, []float64) float64
}

// Approximate evaluates an interpolant at multiple points using multiple
// goroutines.
func Approximate(computer Computer, indices []uint64, surpluses, points []float64,
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
					weight := computer.Compute(indices[k*ni:(k+1)*ni], point)
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
