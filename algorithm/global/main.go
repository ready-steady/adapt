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
	active := []bool{true}
	depths := []uint{0}
	forward := repeatUint(none, 1*ni)
	backward := repeatUint(none, 1*ni)
	progress.Active++

	indices := self.grid.Index(lindices)

	nn := uint(len(indices)) / ni
	nodes := self.grid.Compute(indices)
	counts := []uint{nn}

	values := internal.Invoke(target.Compute, nodes, ni, no, nw)
	progress.Evaluations += nn

	surrogate.push(indices, values)

	lower, upper := updateBounds(nil, nil, values, no)
	scores, errors := updateScores(nil, nil, counts, values, no)
	for {
		target.Monitor(&progress)

		cursor := find(active)

		terminate := true
		δ := threshold(lower, upper, config.AbsTolerance, config.RelTolerance)

	accuracy:
		for _, i := range cursor {
			for j := uint(0); j < no; j++ {
				if errors[i*no+j] > δ[j] {
					terminate = false
					break accuracy
				}
			}
		}
		if terminate {
			break
		}

		min, current := minUint(depths, cursor...)
		max, _ := maxUint(depths)
		if float64(min) > (1.0-config.Adaptivity)*float64(max) {
			_, current = maxFloat64(scores, cursor...)
		}

		active[current] = false
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

			active = append(active, true)
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

		lower, upper = updateBounds(lower, upper, values, no)
		scores, errors = updateScores(scores, errors, counts, surpluses, no)
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

func threshold(lower, upper []float64, absolute, relative float64) []float64 {
	no := uint(len(lower))
	threshold := make([]float64, no)
	for i := uint(0); i < no; i++ {
		threshold[i] = relative * (upper[i] - lower[i])
		if threshold[i] < absolute {
			threshold[i] = absolute
		}
	}
	return threshold
}

func updateBounds(lower, upper []float64, data []float64, no uint) ([]float64, []float64) {
	if lower == nil {
		lower = repeatFloat64(infinity, no)
	}
	if upper == nil {
		upper = repeatFloat64(-infinity, no)
	}
	nn := uint(len(data)) / no
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < no; j++ {
			point := data[i*no+j]
			if lower[j] > point {
				lower[j] = point
			}
			if upper[j] < point {
				upper[j] = point
			}
		}
	}
	return lower, upper
}

func updateScores(scores, errors []float64, counts []uint, surpluses []float64,
	no uint) ([]float64, []float64) {

	offset := uint(0)
	for _, count := range counts {
		score := 0.0
		error := repeatFloat64(-infinity, no)
		for j := uint(0); j < count; j++ {
			for l := uint(0); l < no; l++ {
				Δ := math.Abs(surpluses[(offset+j)*no+l])
				error[l] = math.Max(error[l], Δ)
				score += Δ
			}
		}
		score /= float64(count)
		scores = append(scores, score)
		errors = append(errors, error...)
		offset += count
	}
	return scores, errors
}
