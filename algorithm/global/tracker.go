package global

type Reference map[uint]uint

type Set map[uint]bool

type Tracker struct {
	Active Set
	Length uint

	ni uint

	lmax uint
	imax uint
	rate float64

	lindices []uint64
	norms    []uint64
	scores   []float64

	forward  Reference
	backward Reference
}

func NewTracker(ni, lmax, imax uint, rate float64) *Tracker {
	return &Tracker{
		ni: ni,

		lmax: lmax,
		imax: imax,
		rate: rate,

		forward:  make(Reference),
		backward: make(Reference),
	}
}

func (self *Tracker) Pull() []uint64 {
	if self.Active == nil {
		return self.pullFirst()
	} else {
		return self.pullSubsequent()
	}
}

func (self *Tracker) Push(score float64) {
	self.scores = append(self.scores, score)
}

func (self *Tracker) pullFirst() []uint64 {
	self.Active = make(Set)
	self.Active[0] = true
	self.Length = 1
	self.lindices = make([]uint64, 1*self.ni)
	self.norms = make([]uint64, 1)
	return self.lindices
}

func (self *Tracker) pullSubsequent() (lindices []uint64) {
	ni, nn := self.ni, self.Length
	active, forward, backward := self.Active, self.forward, self.backward

	min, k := minUint64Set(self.norms, active)
	max := maxUint64(self.norms)
	if float64(min) > (1.0-self.rate)*float64(max) {
		_, k = maxFloat64Set(self.scores, active)
	}
	delete(active, k)

	lindex, norm := self.lindices[k*ni:(k+1)*ni], self.norms[k]+1

outer:
	for i := uint(0); i < ni && nn < self.imax; i++ {
		if lindex[i] >= uint64(self.lmax) {
			continue
		}

		newBackward := make(Reference)
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

	lindices = self.lindices[self.Length*ni:]
	self.Length = nn

	return
}
