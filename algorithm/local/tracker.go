package local

type tracker struct {
	ni   uint
	lmin uint
	lmax uint
}

func newTracker(ni uint, config *Config) *tracker {
	return &tracker{
		ni:   ni,
		lmin: config.MinLevel,
		lmax: config.MaxLevel,
	}
}

func (self *tracker) filter(indices []uint64, scores []float64) []uint64 {
	ni, nn := self.ni, uint(len(scores))
	levels := levelize(indices, ni)

	na, ne := uint(0), nn
	for i, j := uint(0), uint(0); i < nn; i++ {
		if levels[i] >= self.lmin && (scores[i] <= 0.0 || levels[i] >= self.lmax) {
			j++
			continue
		}
		if j > na {
			copy(indices[na*ni:], indices[j*ni:ne*ni])
			ne -= j - na
			j = na
		}
		na++
		j++
	}

	return indices[:na*ni]
}
