// Package internal contains types and functions that are shared by the
// interpolation algorithm and are used only internally.
package internal

import (
	"sync"
)

// Basis is a functional basis.
type Basis interface {
	// Compute evaluates the value of a basis function.
	Compute([]uint64, []float64) float64

	// Integrate computes the integral of a basis function.
	Integrate([]uint64) float64
}

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

// Cumulate accumulates weighted integrals of basis functions.
func Cumulate(basis Basis, indices []uint64, surpluses []float64, ni, no uint,
	integral []float64) {

	nn := uint(len(indices)) / ni
	for i := uint(0); i < nn; i++ {
		volume := basis.Integrate(indices[i*ni : (i+1)*ni])
		for j := uint(0); j < no; j++ {
			integral[j] += surpluses[i*no+j] * volume
		}
	}
}
