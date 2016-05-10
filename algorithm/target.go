package algorithm

import (
	"sync"

	"github.com/ready-steady/adapt/algorithm/internal"
)

// Target is a function to be interpolated.
type Target func([]float64, []float64)

// Invoke evaluates a function at multiple points using multiple goroutines.
func Invoke(target Target, points []float64, ni, no uint) []float64 {
	np := uint(len(points)) / ni

	values := make([]float64, np*no)

	jobs := make(chan uint, np)
	group := sync.WaitGroup{}
	group.Add(int(np))

	for i := uint(0); i < internal.Workers; i++ {
		go func() {
			for j := range jobs {
				target(points[j*ni:(j+1)*ni], values[j*no:(j+1)*no])
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
