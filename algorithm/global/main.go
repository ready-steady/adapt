// Package global provides an algorithm for globally adaptive hierarchical
// interpolation.
package global

import (
	"fmt"
	"math"

	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1.0)
)

// Basis is a functional basis.
type Basis interface {
	// Compute evaluates the value of a basis function.
	Compute([]uint64, []float64) float64
}

// Grid is a sparse grid.
type Grid interface {
	// Compute returns the nodes corresponding to a set of indices.
	Compute([]uint64) []float64

	// Index returns the indices of a set of levels.
	Index([]uint8) []uint64
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	config Config
	basis  Basis
	grid   Grid
}

// Progress contains information about the interpolation process.
type Progress struct {
	Level       uint8 // Reached level
	Active      uint  // The number of active indices
	Passive     uint  // The number of passive indices
	Evaluations uint  // The number of function evaluations
}

type cursor map[uint]bool

type reference map[uint]uint

// New creates an interpolator.
func New(grid Grid, basis Basis, config *Config) *Interpolator {
	return &Interpolator{
		config: *config,
		basis:  basis,
		grid:   grid,
	}
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target Target) *Surrogate {
	config := &self.config

	ni, no := target.Dimensions()
	nw := config.Workers

	surrogate := newSurrogate(ni, no)
	progress := newProgress()

	lindices := repeatUint8(0, 1*ni)
	active := make(cursor)
	depths := []uint{0}
	forward := make(reference)
	backward := make(reference)

	active[0] = true
	progress.Active++

	indices := self.grid.Index(lindices)

	nn := uint(len(indices)) / ni
	nodes := self.grid.Compute(indices)
	counts := []uint{nn}

	values := internal.Invoke(target.Compute, nodes, ni, no, nw)
	progress.Evaluations += nn

	surrogate.push(indices, values)

	terminator := newTerminator(no, config)
	terminator.push(values, values, counts)

	scores := assess(target, progress, values, counts, no)
	for !terminator.done(active) {
		target.Monitor(progress)

		min, current := minUint(depths, active)
		max, _ := maxUint(depths)
		if float64(min) > (1.0-config.Adaptivity)*float64(max) {
			_, current = maxFloat64(scores, active)
		}

		delete(active, current)
		progress.Active--
		progress.Passive++

		lindex := lindices[current*ni : (current+1)*ni]

		indices := make([]uint64, 0)
		counts := make([]uint, 0)
		total := progress.Active + progress.Passive

	admissibility:
		for i := uint(0); i < ni && total < config.MaxIndices; i++ {
			if lindex[i] >= config.MaxLevel {
				continue
			}

			newBackward := make(reference)
			for j := uint(0); j < ni; j++ {
				if i == j || lindex[j] == 0 {
					continue
				}
				if l, ok := forward[backward[current*ni+j]*ni+i]; !ok || active[l] {
					continue admissibility
				} else {
					newBackward[j] = l
				}
			}
			newBackward[i] = current
			for j, l := range newBackward {
				forward[l*ni+j] = total
				backward[total*ni+j] = l
			}

			lindices = append(lindices, lindex...)
			lindex := lindices[total*ni:]
			lindex[i]++

			newIndices := self.grid.Index(lindex)
			indices = append(indices, newIndices...)
			counts = append(counts, uint(len(newIndices))/ni)

			active[total] = true
			depths = append(depths, depths[current]+1)

			if lindex[i] > progress.Level {
				progress.Level = lindex[i]
			}

			progress.Active++
			total++
		}

		nn := uint(len(indices)) / ni
		if progress.Evaluations+nn > config.MaxEvaluations {
			break
		}

		nodes := self.grid.Compute(indices)

		values := internal.Invoke(target.Compute, nodes, ni, no, nw)
		progress.Evaluations += nn

		surpluses := internal.Subtract(values, internal.Approximate(self.basis,
			surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.push(indices, surpluses)
		scores = append(scores, assess(target, progress, surpluses, counts, no)...)

		terminator.push(values, surpluses, counts)
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}

func newProgress() *Progress {
	return &Progress{}
}

// String returns a human-friendly representation.
func (self *Progress) String() string {
	phantom := struct {
		level       uint8
		active      uint
		passive     uint
		evaluations uint
	}{
		level:       self.Level,
		active:      self.Active,
		passive:     self.Passive,
		evaluations: self.Evaluations,
	}
	return fmt.Sprintf("%+v", phantom)
}
