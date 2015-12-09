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
	none     = ^uint(0)
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
	progress := Progress{}

	lindices := repeatUint8(0, 1*ni)
	active := make(cursor)
	depths := []uint{0}
	forward := repeatUint(none, 1*ni)
	backward := repeatUint(none, 1*ni)

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

	scores := updateScores(nil, counts, values, no)
	for {
		target.Monitor(&progress)

		if terminator.check(active) {
			break
		}

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

			newBackward := repeatUint(none, ni)
			newBackward[i] = current
			for j := uint(0); j < ni; j++ {
				if i == j || lindex[j] == 0 {
					continue
				}
				l := forward[backward[current*ni+j]*ni+i]
				if l == none || active[l] {
					continue admissibility
				}
				newBackward[j] = l
			}

			lindices = append(lindices, lindex...)
			lindex := lindices[total*ni:]
			lindex[i]++

			newIndices := self.grid.Index(lindex)
			indices = append(indices, newIndices...)
			counts = append(counts, uint(len(newIndices))/ni)

			for j := uint(0); j < ni; j++ {
				if newBackward[j] != none {
					forward[newBackward[j]*ni+j] = total
				}
			}

			active[total] = true
			depths = append(depths, depths[current]+1)
			forward = append(forward, repeatUint(none, ni)...)
			backward = append(backward, newBackward...)

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

		terminator.push(values, surpluses, counts)
		scores = updateScores(scores, counts, surpluses, no)
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
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

func updateScores(scores []float64, counts []uint, surpluses []float64, no uint) []float64 {
	for _, count := range counts {
		scores = append(scores, score(surpluses[:count*no], no))
		surpluses = surpluses[count*no:]
	}
	return scores
}

func score(surpluses []float64, no uint) float64 {
	score := 0.0
	for _, value := range surpluses {
		if value < 0.0 {
			value = -value
		}
		score += value
	}
	return score / float64(uint(len(surpluses))/no)
}
