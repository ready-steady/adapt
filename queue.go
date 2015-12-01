package adapt

import (
	"math"

	"github.com/ready-steady/sort"
)

var (
	infinity = math.Inf(1)
)

type queue struct {
	ni uint
	no uint
	nn uint

	lmin uint
	lmax uint
	rate float64

	Indices []uint64
	Nodes   []float64
	Values  []float64
	Scores  []float64
}

func newQueue(ni, no uint, config *Config) *queue {
	queue := &queue{
		ni:   ni,
		no:   no,
		lmin: config.MinLevel,
		lmax: config.MaxLevel,
		rate: config.Rate,
	}
	queue.empty()
	return queue
}

func (self *queue) compress(from uint) {
	ni, no, nn := self.ni, self.no, self.nn-from

	indices := self.Indices[from*ni:]
	nodes := self.Nodes[from*ni:]
	values := self.Values[from*no:]
	scores := self.Scores[from:]

	na, ne := uint(0), nn
	for i, j := uint(0), uint(0); i < nn; i++ {
		if scores[j] <= 0.0 {
			j++
			continue
		}
		if j > na {
			copy(indices[na*ni:], indices[j*ni:ne*ni])
			copy(nodes[na*ni:], nodes[j*ni:ne*ni])
			copy(values[na*no:], values[j*no:ne*no])
			copy(scores[na:], scores[j:ne])
			ne -= j - na
			j = na
		}
		na++
		j++
	}

	nn = from + na
	self.nn = nn

	self.Indices = self.Indices[:nn*ni]
	self.Nodes = self.Nodes[:nn*ni]
	self.Values = self.Values[:nn*no]
	self.Scores = self.Scores[:nn]
}

func (self *queue) empty() {
	self.nn = 0
	self.Indices = []uint64{}
	self.Nodes = []float64{}
	self.Values = []float64{}
	self.Scores = []float64{}
}

func (self *queue) pull() []uint64 {
	ni, nn := self.ni, uint(math.Ceil(self.rate*float64(self.nn)))

	if nn == self.nn {
		indices := self.Indices
		self.empty()
		return indices
	}

	scores := make([]float64, self.nn)
	copy(scores, self.Scores)
	order, _ := sort.Quick(scores)

	indices := make([]uint64, nn*ni)
	for i, j := uint(0), self.nn-nn; i < nn; i, j = i+1, j+1 {
		k := order[j]
		copy(indices[i*ni:(i+1)*ni], self.Indices[k*ni:(k+1)*ni])
		self.Scores[k] = 0.0
	}
	self.compress(0)

	return indices
}

func (self *queue) push(indices []uint64, nodes, values, scores []float64) {
	nn := self.nn

	self.Indices = append(self.Indices, indices...)
	self.Nodes = append(self.Nodes, nodes...)
	self.Values = append(self.Values, values...)
	self.Scores = append(self.Scores, scores...)

	self.nn += uint(len(scores))

	for i, level := range levelize(indices, self.ni) {
		if level < self.lmin {
			self.Scores[nn+uint(i)] = infinity
		} else if level >= self.lmax {
			self.Scores[nn+uint(i)] = 0.0
		}
	}

	self.compress(nn)
}

func (self *queue) update(scores []float64) {
	copy(self.Scores, scores)
	self.compress(0)
}
