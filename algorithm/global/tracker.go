package global

type reference map[uint]uint

type tracker struct {
	ni uint
	nn uint

	lmax uint8
	imax uint

	lindices []uint8
	forward  reference
	backward reference
}

func newTracker(ni uint, config *Config) *tracker {
	return &tracker{
		ni: ni,

		lmax: config.MaxLevel,
		imax: config.MaxIndices,

		forward:  make(reference),
		backward: make(reference),
	}
}

func (self *tracker) pull(position uint, active cursor) (lindices []uint8) {
	ni, nn := self.ni, self.nn
	forward, backward := self.forward, self.backward

	lindex := self.lindices[position*ni : (position+1)*ni]

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
			if l, ok := forward[backward[position*ni+j]*ni+i]; !ok || active[l] {
				continue outer
			} else {
				newBackward[j] = l
			}
		}
		newBackward[i] = position
		for j, l := range newBackward {
			forward[l*ni+j] = nn
			backward[nn*ni+j] = l
		}

		self.lindices = append(self.lindices, lindex...)
		self.lindices[nn*ni+i]++
		nn++
	}

	lindices = self.lindices[self.nn*ni:]
	self.nn = nn

	return
}

func (self *tracker) push(lindices []uint8) {
	self.nn += uint(len(lindices)) / self.ni
	self.lindices = append(self.lindices, lindices...)
}
