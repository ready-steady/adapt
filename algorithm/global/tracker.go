package global

type reference map[uint]uint

type tracker struct {
	ni uint
	nn uint

	lmax uint
	imax uint
	rate float64

	lindices []uint64
	norms    []uint64
	scores   []float64

	active   Set
	forward  reference
	backward reference
}

func newTracker(ni uint, config *Config) *tracker {
	return &tracker{
		ni: ni,

		lmax: config.MaxLevel,
		imax: config.MaxIndices,
		rate: config.AdaptivityRate,

		forward:  make(reference),
		backward: make(reference),
	}
}

func (self *tracker) pull() []uint64 {
	if self.lindices == nil {
		return self.pullFirst()
	} else {
		return self.pullSubsequent()
	}
}

func (self *tracker) pullFirst() []uint64 {
	self.lindices = make([]uint64, 1*self.ni)
	self.norms = make([]uint64, 1)
	self.active = make(Set)
	self.active[0] = true
	self.nn = 1
	return self.lindices
}

func (self *tracker) pullSubsequent() (lindices []uint64) {
	ni, nn := self.ni, self.nn
	active, forward, backward := self.active, self.forward, self.backward

	min, k := minUint64Set(self.norms, self.active)
	max := maxUint64(self.norms)
	if float64(min) > (1.0-self.rate)*float64(max) {
		_, k = maxFloat64Set(self.scores, self.active)
	}
	delete(active, k)

	lindex, norm := self.lindices[k*ni:(k+1)*ni], self.norms[k]+1

outer:
	for i := uint(0); i < ni && nn < self.imax; i++ {
		if lindex[i] >= uint64(self.lmax) {
			continue
		}

		newBackward := make(reference)
		for j := uint(0); j < ni; j++ {
			if i == j || lindex[j] == 0 {
				continue
			}
			if l, ok := forward[backward[k*ni+j]*ni+i]; !ok || active[l] {
				continue outer
			} else {
				newBackward[j] = l
			}
		}
		newBackward[i] = k
		for j, l := range newBackward {
			forward[l*ni+j] = nn
			backward[nn*ni+j] = l
		}

		self.lindices = append(self.lindices, lindex...)
		self.lindices[nn*ni+i]++
		self.norms = append(self.norms, norm)

		active[nn] = true

		nn++
	}

	lindices = self.lindices[self.nn*ni:]
	self.nn = nn

	return
}

func (self *tracker) push(score float64) {
	self.scores = append(self.scores, score)
}
