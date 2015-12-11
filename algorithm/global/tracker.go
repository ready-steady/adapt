package global

type reference map[uint]uint

type tracker struct {
	ni uint
	nn uint

	lmax       uint8
	imax       uint
	adaptivity float64

	lindices []uint8
	norms    []uint
	scores   []float64

	active   cursor
	forward  reference
	backward reference
}

func newTracker(ni uint, config *Config) *tracker {
	return &tracker{
		ni: ni,

		lmax:       config.MaxLevel,
		imax:       config.MaxIndices,
		adaptivity: config.Adaptivity,

		forward:  make(reference),
		backward: make(reference),
	}
}

func (self *tracker) pull() []uint8 {
	if self.lindices == nil {
		return self.pullFirst()
	} else {
		return self.pullSubsequent()
	}
}

func (self *tracker) pullFirst() []uint8 {
	self.lindices = make([]uint8, 1*self.ni)
	self.norms = make([]uint, 1)
	self.active = make(cursor)
	self.active[0] = true
	self.nn = 1
	return self.lindices
}

func (self *tracker) pullSubsequent() (lindices []uint8) {
	ni, nn := self.ni, self.nn
	active, forward, backward := self.active, self.forward, self.backward

	min, k := minUint(self.norms, self.active)
	max := maxUint(self.norms)
	if float64(min) > (1.0-self.adaptivity)*float64(max) {
		_, k = maxFloat64(self.scores, self.active)
	}
	delete(active, k)

	lindex, norm := self.lindices[k*ni:(k+1)*ni], self.norms[k]+1

outer:
	for i := uint(0); i < ni && nn < self.imax; i++ {
		if lindex[i] >= self.lmax {
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

func (self *tracker) push(scores []float64) {
	self.scores = append(self.scores, scores...)
}
