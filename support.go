package adapt

import (
	"sync"
)

func approximate(basis Basis, indices []uint64, surpluses, points []float64,
	ni, no, nw uint) []float64 {

	nn, np := uint(len(indices))/ni, uint(len(points))/ni

	values := make([]float64, np*no)

	jobs := make(chan uint, np)
	group := sync.WaitGroup{}
	group.Add(int(np))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				point := points[j*ni : (j+1)*ni]
				value := values[j*no : (j+1)*no]

				for k := uint(0); k < nn; k++ {
					weight := basis.Compute(indices[k*ni:(k+1)*ni], point)
					if weight == 0 {
						continue
					}
					for l := uint(0); l < no; l++ {
						value[l] += weight * surpluses[k*no+l]
					}
				}

				group.Done()
			}
		}()
	}

	for i := uint(0); i < np; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}

func balance(grid Grid, history *hash, indices []uint64) []uint64 {
	neighbors := make([]uint64, 0)

	for {
		indices = socialize(grid, history, indices)

		if len(indices) == 0 {
			break
		}

		neighbors = append(neighbors, indices...)
	}

	return neighbors
}

func compact(indices []uint64, surpluses, scores []float64,
	ni, no, nn uint) ([]uint64, []float64, []float64) {

	na, ne := uint(0), nn
	for i, j := uint(0), uint(0); i < nn; i++ {
		if scores[j] < 0 {
			j++
			continue
		}

		if j > na {
			copy(indices[na*ni:], indices[j*ni:ne*ni])
			copy(surpluses[na*no:], surpluses[j*no:ne*no])
			copy(scores[na:], scores[j:ne])
			ne -= j - na
			j = na
		}

		na++
		j++
	}

	return indices[:na*ni], surpluses[:na*no], scores[:na]
}

func cumulate(basis Basis, indices []uint64, surpluses []float64, ni, no, nn uint,
	integral []float64) {

	for i := uint(0); i < nn; i++ {
		volume := basis.Integrate(indices[i*ni : (i+1)*ni])
		for j := uint(0); j < no; j++ {
			integral[j] += surpluses[i*no+j] * volume
		}
	}
}

func invoke(compute func([]float64, []float64), nodes []float64, ni, no, nw uint) []float64 {
	nn := uint(len(nodes)) / ni

	values := make([]float64, nn*no)

	jobs := make(chan uint, nn)
	group := sync.WaitGroup{}
	group.Add(int(nn))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				compute(nodes[j*ni:(j+1)*ni], values[j*no:(j+1)*no])
				group.Done()
			}
		}()
	}

	for i := uint(0); i < nn; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}

func measure(basis Basis, indices []uint64, ni uint) []float64 {
	nn := uint(len(indices)) / ni

	volumes := make([]float64, nn)
	for i := uint(0); i < nn; i++ {
		volumes[i] = basis.Integrate(indices[i*ni : (i+1)*ni])
	}

	return volumes
}

func socialize(grid Grid, history *hash, indices []uint64) []uint64 {
	ni := history.ni
	nn := uint(len(indices)) / ni

	siblings := make([]uint64, 0, ni)
	for i := uint(0); i < nn; i++ {
		index := indices[i*ni : (i+1)*ni]

		for j := uint(0); j < ni; j++ {
			pair := index[j]

			grid.Parent(index, j)
			if !history.find(index) {
				index[j] = pair
				continue
			}
			index[j] = pair

			grid.Sibling(index, j)
			if !history.find(index) {
				history.push(index)
				siblings = append(siblings, index...)
			}
			index[j] = pair
		}
	}

	return siblings
}

func subtract(minuend, subtrahend []float64) []float64 {
	difference := make([]float64, len(minuend))
	for i := range minuend {
		difference[i] = minuend[i] - subtrahend[i]
	}
	return difference
}
