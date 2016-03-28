// Package internal contains code shared by the neighbor packages.
package internal

import (
	"runtime"
	"sync"
)

var (
	// Workers is the number of goroutines used for interpolation.
	Workers = uint(runtime.GOMAXPROCS(0))
)

// BasisComputer returns the value of a basis function.
type BasisComputer interface {
	Compute([]uint64, []float64) float64
}

// BasisIntegrator returns the integral of a basis function.
type BasisIntegrator interface {
	Integrate([]uint64) float64
}

// GridComputer returns the coordinates of a grid node.
type GridComputer interface {
	Compute([]uint64) []float64
}

// GridIndexer returns the nodal indices of a level index.
type GridIndexer interface {
	Index([]uint64) []uint64
}

// GridRefiner returns the child nodal indices of a set of nodal indices.
type GridRefiner interface {
	Refine([]uint64) []uint64
}

// GridRefinerToward returns the child nodal indices of a set of nodal indices
// with respect to a particular dimension.
type GridRefinerToward interface {
	RefineToward([]uint64, uint) []uint64
}

// Approximate evaluates an interpolant at multiple points using multiple
// goroutines.
func Approximate(computer BasisComputer, indices []uint64, surpluses, points []float64,
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

// Index returns the nodal indices of a set of level indices.
func Index(indexer GridIndexer, lindices []uint64, ni uint) ([]uint64, []uint) {
	nn := uint(len(lindices)) / ni
	indices, counts := []uint64(nil), make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		newIndices := indexer.Index(lindices[i*ni : (i+1)*ni])
		indices = append(indices, newIndices...)
		counts[i] = uint(len(newIndices)) / ni
	}
	return indices, counts
}

// Measure computes the integrals of a set of basis functions.
func Measure(integrator BasisIntegrator, indices []uint64, ni uint) []float64 {
	nn := uint(len(indices)) / ni
	volumes := make([]float64, nn)
	for i := uint(0); i < nn; i++ {
		volumes[i] = integrator.Integrate(indices[i*ni : (i+1)*ni])
	}
	return volumes
}
