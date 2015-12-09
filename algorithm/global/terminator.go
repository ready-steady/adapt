package global

type terminator struct {
	no       uint
	absolute float64
	relative float64

	errors []float64
	lower  []float64
	upper  []float64
}

func newTerminator(no uint, config *Config) *terminator {
	return &terminator{
		no:       no,
		absolute: config.AbsTolerance,
		relative: config.RelTolerance,

		lower: repeatFloat64(infinity, no),
		upper: repeatFloat64(-infinity, no),
	}
}

func (self *terminator) done(cursor cursor) bool {
	no, errors := self.no, self.errors
	δ := threshold(self.lower, self.upper, self.absolute, self.relative)
	for i := range cursor {
		for j := uint(0); j < no; j++ {
			if errors[i*no+j] > δ[j] {
				return false
			}
		}
	}
	return true
}

func (self *terminator) push(values, surpluses []float64, counts []uint) {
	no := self.no
	for i, point := range values {
		j := uint(i) % no
		if self.lower[j] > point {
			self.lower[j] = point
		}
		if self.upper[j] < point {
			self.upper[j] = point
		}
	}
	for _, count := range counts {
		self.errors = append(self.errors, error(surpluses[:count*no], no)...)
		surpluses = surpluses[count*no:]
	}
}

func error(surpluses []float64, no uint) []float64 {
	error := repeatFloat64(-infinity, no)
	for i, value := range surpluses {
		j := uint(i) % no
		if value < 0.0 {
			value = -value
		}
		if value > error[j] {
			error[j] = value
		}
	}
	return error
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
