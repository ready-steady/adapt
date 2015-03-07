// Package adhier provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adhier

import (
	"math"
	"runtime"
	"sync"
)

// Grid is the interface that an sparse grid should satisfy in order to be used
// in the algorithm.
type Grid interface {
	ComputeNodes(indices []uint64) []float64
	ComputeChildren(indices []uint64) []uint64
}

// Basis is the interface that a functional basis should satisfy in order to be
// used in the algorithm.
type Basis interface {
	Evaluate(index []uint64, point []float64) float64
}

// Interpolator represents a particular instantiation of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config *Config

	ni uint
	no uint
	nw uint
}

// New creates an instance of the algorithm for the given configuration.
func New(grid Grid, basis Basis, config *Config) (*Interpolator, error) {
	if err := config.verify(); err != nil {
		return nil, err
	}

	nw := config.Workers
	if nw == 0 {
		nw = uint(runtime.GOMAXPROCS(0))
	}

	interpolator := &Interpolator{
		grid:   grid,
		basis:  basis,
		config: config,

		ni: config.Inputs,
		no: config.Outputs,
		nw: nw,
	}

	return interpolator, nil
}

// Compute takes a target function and produces an interpolant for it. The
// interpolant can then be fed to Evaluate for approximating the target function
// at arbitrary points.
//
// The second argument of Compute is an optional function that can be used for
// monitoring the progress of interpolation. The progress function is called
// once for each level before evaluating the target function at the nodes of
// that level. The signature of the progress function is func(uint, uint, uint)
// where the arguments are the current level, number of active nodes, and total
// number of nodes, respectively.
func (self *Interpolator) Compute(target func([]float64, []float64, []uint64),
	arguments ...interface{}) *Surrogate {

	var progress func(uint, uint, uint)
	if len(arguments) > 0 {
		progress = arguments[0].(func(uint, uint, uint))
	}

	ni, no := self.ni, self.no
	config := self.config

	surrogate := new(Surrogate)
	surrogate.initialize(ni, no)

	// Level 0 is assumed to have only one node, and the order of that node is
	// assumed to be zero.
	level := uint(0)

	na := uint(1) // active
	np := uint(0) // passive

	indices := make([]uint64, na*ni)

	var i, j, k, l uint
	var nodes, values, approximations []float64

	min := make([]float64, no)
	max := make([]float64, no)

	min[0], max[0] = math.Inf(1), math.Inf(-1)
	for i = 1; i < no; i++ {
		min[i], max[i] = min[0], max[0]
	}

	for {
		if progress != nil {
			progress(level, na, np+na)
		}

		surrogate.resize(np + na)
		copy(surrogate.Indices[np*ni:], indices)

		nodes = self.grid.ComputeNodes(indices)

		// NOTE: Assuming that target might have some logic based on the indices
		// passed to it (for instance, caching), the indices variable should not
		// be used here as it gets modified later on.
		values = self.invoke(target, nodes, surrogate.Indices[np*ni:(np+na)*ni])

		// Compute the surpluses corresponding to the active nodes.
		if level == 0 {
			// The surrogate does not have any nodes yet.
			copy(surrogate.Surpluses, values)
			goto refineLevel
		}

		approximations = self.approximate(surrogate.Indices[:np*ni],
			surrogate.Surpluses[:np*no], nodes)
		for i, k = 0, np*no; i < na; i++ {
			for j = 0; j < no; j++ {
				surrogate.Surpluses[k] = values[i*no+j] - approximations[i*no+j]
				k++
			}
		}

	refineLevel:
		if level >= config.MaxLevel || (np+na) >= config.MaxNodes {
			break
		}

		// Keep track of the maximal and minimal values of the function.
		for i, k = 0, 0; i < na; i++ {
			for j = 0; j < no; j++ {
				if values[k] < min[j] {
					min[j] = values[k]
				}
				if values[k] > max[j] {
					max[j] = values[k]
				}
				k++
			}
		}

		if level < config.MinLevel {
			goto updateIndices
		}

		k, l = 0, 0

		for i = 0; i < na; i++ {
			refine := false

			for j = 0; j < no; j++ {
				absError := surrogate.Surpluses[(np+i)*no+j]
				if absError < 0 {
					absError = -absError
				}

				if absError > config.AbsError {
					refine = true
					break
				}

				relError := absError / (max[j] - min[j])

				if relError > config.RelError {
					refine = true
					break
				}
			}

			if !refine {
				l += ni
				continue
			}

			if k != l {
				// Shift everything, assuming a lot of refinements.
				copy(indices[k:], indices[l:])
				l = k
			}

			k += ni
			l += ni
		}

		indices = indices[:k]

	updateIndices:
		indices = self.grid.ComputeChildren(indices)

		np += na
		na = uint(len(indices)) / ni

		if Δ := int32(np+na) - int32(config.MaxNodes); Δ > 0 {
			na -= uint(Δ)
			indices = indices[:na*ni]
		}
		if na == 0 {
			break
		}

		level++
	}

	surrogate.finalize(level, np+na)
	return surrogate
}

// Evaluate takes a surrogate produced by Compute and evaluates it at a set of
// points.
func (self *Interpolator) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	return self.approximate(surrogate.Indices, surrogate.Surpluses, points)
}

func (self *Interpolator) approximate(indices []uint64, surpluses, points []float64) []float64 {
	ni, no, nw := self.ni, self.no, self.nw
	nn := uint(len(indices)) / ni
	np := uint(len(points)) / ni

	basis := self.basis

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
					weight := basis.Evaluate(indices[k*ni:(k+1)*ni], point)
					if weight == 0 {
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

func (self *Interpolator) invoke(target func([]float64, []float64, []uint64),
	nodes []float64, indices []uint64) []float64 {

	ni, no, nw := self.ni, self.no, self.nw
	nn := uint(len(nodes)) / ni

	values := make([]float64, nn*no)

	jobs := make(chan uint, nn)
	group := sync.WaitGroup{}
	group.Add(int(nn))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				target(nodes[j*ni:(j+1)*ni], values[j*no:(j+1)*no], indices[j*ni:(j+1)*ni])
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
