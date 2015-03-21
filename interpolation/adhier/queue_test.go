package adhier

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestQueuePushPull(t *testing.T) {
	config := NewConfig()
	config.Rate = 0.5

	queue := newQueue(1, config)

	assert.Equal(queue.pull(), []uint64{0}, t)

	queue.push([]uint64{1, 2, 3, 4, 5, 6}, []float64{2, 0, 4, 3, 1, 5})

	assert.Equal(queue.pull(), []uint64{6, 3, 4}, t)
	assert.Equal(queue.pull(), []uint64{1}, t)
	assert.Equal(queue.pull(), []uint64{5}, t)
}
