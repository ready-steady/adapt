package global

type selector struct {
	ni         uint
	adaptivity float64

	depths []uint
	scores []float64
}

func newSelector(ni uint, config *Config) *selector {
	return &selector{
		ni:         ni,
		adaptivity: config.Adaptivity,
	}
}

func (self *selector) pull(cursor cursor) (uint, uint) {
	min, position := minUint(self.depths, cursor)
	max, _ := maxUint(self.depths)
	if float64(min) > (1.0-self.adaptivity)*float64(max) {
		_, position = maxFloat64(self.scores, cursor)
	}
	return position, self.depths[position]
}

func (self *selector) push(scores []float64, depth uint) {
	self.depths = append(self.depths, repeatUint(depth, uint(len(scores)))...)
	self.scores = append(self.scores, scores...)
}
