// Package internal contains code shared by the neighbor packages.
package internal

import (
	"runtime"
	"sync"

	"github.com/ready-steady/adapt/basis"
	"github.com/ready-steady/adapt/grid"
)

var (
	// Workers is the number of goroutines used for interpolation.
	Workers = uint(runtime.GOMAXPROCS(0))
)

// Approximate evaluates an interpolant at multiple points using multiple
// goroutines.
func Approximate(computer basis.Computer, indices []uint64, surpluses, points []float64,
	ni, no uint) []float64 {

	nn := uint(len(indices)) / ni
	np := uint(len(points)) / ni
	values := make([]float64, np*no)

	jobs := make(chan uint, np)
	group := sync.WaitGroup{}
	group.Add(int(np))

	for i := uint(0); i < Workers; i++ {
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

// Index returns the nodal indices of a set of level indices.
func Index(indexer grid.Indexer, ildices []uint64, ni uint) ([]uint64, []uint) {
	nn := uint(len(ildices)) / ni
	indices, counts := []uint64(nil), make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		more := indexer.Index(ildices[i*ni : (i+1)*ni])
		indices = append(indices, more...)
		counts[i] = uint(len(more)) / ni
	}
	return indices, counts
}

// Measure computes the integrals of a set of basis functions.
func Measure(integrator basis.Integrator, indices []uint64, ni uint) []float64 {
	nn := uint(len(indices)) / ni
	volumes := make([]float64, nn)
	for i := uint(0); i < nn; i++ {
		volumes[i] = integrator.Integrate(indices[i*ni : (i+1)*ni])
	}
	return volumes
}
