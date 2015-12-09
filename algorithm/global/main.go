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

	indices := make([]uint64, ni)
	nodes, counts := self.grid.Compute(indices), []uint{1}
	values := internal.Invoke(target.Compute, nodes, ni, no, nw)

	surrogate.push(indices, values)
	progress.Evaluations++

	lindices := repeatUint8(0, ni)
	active := []bool{true}
	depths := []uint{0}
	forward := repeatUint(none, 1*ni)
	backward := repeatUint(none, 1*ni)
	progress.Active++

	lower, upper := updateBounds(nil, nil, values, no)
	scores, errors := updateScores(nil, nil, counts, values, no)
	for {
		target.Monitor(&progress)

		cursor := find(active)

		terminate := true
		δ := threshold(lower, upper, config.AbsTolerance, config.RelTolerance)

	accuracyCheck:
		for _, i := range cursor {
			for j := uint(0); j < no; j++ {
				if errors[i*no+j] > δ[j] {
					terminate = false
					break accuracyCheck
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

		cursor = make([]uint, 0, ni)
		newBackward := repeatUint(none, ni*ni)

	admissibilityCheck:
		for i := uint(0); i < ni; i++ {
			if lindex[i] >= config.MaxLevel {
				continue
			}
			newBackward[i*ni+i] = current
			for j := uint(0); j < ni; j++ {
				if i == j || lindex[j] == 0 {
					continue
				}
				l := forward[backward[current*ni+j]*ni+i]
				if l == none || active[l] {
					continue admissibilityCheck
				}
				newBackward[i*ni+j] = l
			}
			cursor = append(cursor, i)
		}

		na := uint(len(cursor))
		if na == 0 {
			continue
		}

		nt := progress.Active + progress.Passive
		if nt+na > config.MaxIndices {
			break
		}

		indices := make([]uint64, 0)
		counts := make([]uint, 0, na)
		for i := uint(0); i < na; i++ {
			j := cursor[i]

			level := lindex[j] + 1
			if level > progress.Level {
				progress.Level = level
			}

			lindices = append(lindices, lindex...)
			lindices[(nt+i)*ni+j] = level

			newIndices := self.grid.Index(lindices[(nt+i)*ni:])
			indices = append(indices, newIndices...)
			counts = append(counts, uint(len(newIndices))/ni)

			active = append(active, true)
			depths = append(depths, depths[current]+1)

			for l := uint(0); l < ni; l++ {
				if newBackward[j*ni+l] == none {
					continue
				}
				forward[newBackward[j*ni+l]*ni+l] = nt + i
			}
			forward = appendUint(forward, none, ni)
			backward = append(backward, newBackward[j*ni:(j+1)*ni]...)

			progress.Active++
		}

		nn := uint(len(indices)) / ni
		if progress.Evaluations+nn > config.MaxEvaluations {
			break
		}

		nodes := self.grid.Compute(indices)
		values := internal.Invoke(target.Compute, nodes, ni, no, nw)
		surpluses := internal.Subtract(values, internal.Approximate(self.basis,
			surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.push(indices, surpluses)
		progress.Evaluations += nn

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
