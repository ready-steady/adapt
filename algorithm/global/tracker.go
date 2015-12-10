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

func (self *tracker) pull(active cursor) (lindices []uint8) {
	ni, nn := self.ni, self.nn
	forward, backward := self.forward, self.backward

	k := self.choose(active)
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

		nn++
	}

	lindices = self.lindices[self.nn*ni:]
	self.nn = nn

	return
}

func (self *tracker) push(lindices []uint8, scores []float64) {
	ni := self.ni
	nn := uint(len(lindices)) / ni

	self.nn += nn
	self.lindices = append(self.lindices, lindices...)
	self.scores = append(self.scores, scores...)
	for i := uint(0); i < nn; i++ {
		norm := uint(0)
		for j := uint(0); j < ni; j++ {
			norm += uint(lindices[i*ni+j])
		}
		self.norms = append(self.norms, norm)
	}
}

func (self *tracker) choose(active cursor) uint {
	min, k := minUint(self.norms, active)
	max := maxUint(self.norms)
	if float64(min) > (1.0-self.adaptivity)*float64(max) {
		_, k = maxFloat64(self.scores, active)
	}
	return k
}
