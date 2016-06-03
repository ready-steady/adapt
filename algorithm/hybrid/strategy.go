package hybrid

import (
	"github.com/ready-steady/adapt/algorithm"
	"github.com/ready-steady/adapt/algorithm/internal"
	"github.com/ready-steady/adapt/grid"
)

// Strategy is a basic strategy.
type Strategy struct {
	ni uint
	no uint

	guide Guide

	lmin uint
	lmax uint
	εs   float64

	active    *internal.Active
	threshold *internal.Threshold
	hash      *internal.Hash
	unique    *internal.Unique

	priority []float64
	accuracy []float64

	scopes [][]uint
	scores []float64

	lndexer map[string]uint
	indexer map[string]uint
}

// Guide is a grid-refinement tool of a basic strategy.
type Guide interface {
	grid.Indexer
	grid.RefinerToward
}

// NewStrategy creates a basic strategy.
func NewStrategy(inputs, outputs uint, guide Guide, minLevel, maxLevel uint,
	absoluteError, relativeError, scoreError float64) *Strategy {

	return &Strategy{
		ni: inputs,
		no: outputs,

		guide: guide,

		lmin: minLevel,
		lmax: maxLevel,
		εs:   scoreError,

		active:    internal.NewActive(inputs),
		threshold: internal.NewThreshold(outputs, absoluteError, relativeError),
		hash:      internal.NewHash(inputs),
		unique:    internal.NewUnique(inputs),

		lndexer: make(map[string]uint),
		indexer: make(map[string]uint),
	}
}

func (self *Strategy) First(surrogate *algorithm.Surrogate) *algorithm.State {
	return self.initiate(self.active.First(), surrogate)
}

func (self *Strategy) Next(state *algorithm.State,
	surrogate *algorithm.Surrogate) *algorithm.State {

	exclude := make(map[uint]bool)
	for {
		self.consume(state)
		if self.threshold.Check(self.accuracy, self.active.Positions) {
			return nil
		}
		k := internal.Choose(self.priority, self.active.Positions, exclude)
		if k == internal.None {
			return nil
		}
		lndices := self.active.Next(k)
		if len(lndices) > 0 {
			self.active.Drop(k)
		} else {
			exclude[k] = true
		}
		state = self.initiate(lndices, surrogate)
		if len(state.Indices) > 0 {
			return state
		}
	}
}

func (self *Strategy) Score(element *algorithm.Element) float64 {
	return internal.MaxAbsolute(element.Surplus) * element.Volume
}

func (self *Strategy) consume(state *algorithm.State) {
	ni, no := self.ni, self.no
	np := uint(len(self.priority))
	na := uint(len(self.accuracy))
	ns := uint(len(self.scores))
	nn := uint(len(state.Counts))

	self.priority = append(self.priority, make([]float64, nn)...)
	priority := self.priority[np:]

	self.accuracy = append(self.accuracy, make([]float64, nn*no)...)
	accuracy := self.accuracy[na:]

	self.scopes = append(self.scopes, make([][]uint, nn)...)
	scopes := self.scopes[np:]

	self.scores = append(self.scores, make([]float64, len(state.Scores))...)
	scores := self.scores[ns:]

	groups := state.Data.([][]uint64)
	levels := internal.Levelize(state.Lndices, ni)

	for i, o := uint(0), uint(0); i < nn; i++ {
		count := state.Counts[i]
		if levels[i] < uint64(self.lmin) {
			internal.Set(accuracy[i*no:(i+1)*no], internal.Infinity)
			internal.Set(scores[o:(o+count)], internal.Infinity)
		} else if levels[i] < uint64(self.lmax) {
			self.threshold.Compress(accuracy[i*no:(i+1)*no],
				state.Surpluses[o*no:(o+count)*no])
			copy(scores[o:(o+count)], state.Scores[o:(o+count)])
		}
		for j := uint(0); j < count; j++ {
			index := state.Indices[(o+j)*ni : (o+j+1)*ni]
			self.indexer[self.hash.Key(index)] = ns + o + j
		}
		o += count
	}

	for i := uint(0); i < nn; i++ {
		count := uint(len(groups[i])) / ni
		scope := make([]uint, count)
		for j := uint(0); j < count; j++ {
			index := groups[i][j*ni : (j+1)*ni]
			k, ok := self.indexer[self.hash.Key(index)]
			if !ok {
				panic("something went wrong")
			}
			scope[j] = k
		}
		scopes[i] = scope
		if levels[i] < uint64(self.lmin) {
			priority[i] = internal.Infinity
		} else if levels[i] < uint64(self.lmax) {
			for _, j := range scope {
				priority[i] += self.scores[j]
			}
			priority[i] /= float64(count)
		}
		lndex := state.Lndices[i*ni : (i+1)*ni]
		self.lndexer[self.hash.Key(lndex)] = np + i
	}

	self.threshold.Update(state.Values)
}

func (self *Strategy) index(lndices []uint64, surrogate *algorithm.Surrogate) [][]uint64 {
	ni := self.ni
	nn := uint(len(lndices)) / ni
	groups := make([][]uint64, nn)
	for i := uint(0); i < nn; i++ {
		root, lndex := true, lndices[i*ni:(i+1)*ni]
		for j := uint(0); j < ni; j++ {
			level := lndex[j]
			if level == 0 {
				continue
			}
			lndex[j] = level - 1
			k, ok := self.lndexer[self.hash.Key(lndex)]
			lndex[j] = level
			if !ok {
				panic("something went wrong")
			}
			for _, l := range self.scopes[k] {
				if self.scores[l] >= self.εs {
					index := surrogate.Indices[l*ni : (l+1)*ni]
					groups[i] = append(groups[i], self.guide.RefineToward(index, j)...)
				}
			}
			root = false
		}
		if root {
			groups[i] = append(groups[i], self.guide.Index(lndex)...)
		}
	}
	return groups
}

func (self *Strategy) initiate(lndices []uint64,
	surrogate *algorithm.Surrogate) (state *algorithm.State) {

	groups := self.index(lndices, surrogate)
	nn := uint(len(groups))
	state = &algorithm.State{
		Lndices: lndices,
		Counts:  make([]uint, nn),
		Data:    groups,
	}
	for i := uint(0); i < nn; i++ {
		indices := self.unique.Distil(groups[i])
		state.Indices = append(state.Indices, indices...)
		state.Counts[i] = uint(len(indices)) / self.ni
	}
	return
}
