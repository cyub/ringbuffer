package ringbuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSPMCRingBuffer(t *testing.T) {
	cap := 16
	spsc := NewSpmcRingBuffer(cap)

	assert.Equal(t, spsc.Length(), 0)
	assert.Equal(t, spsc.Capacity(), cap)

	for i := 0; i < cap*2; i++ {
		err := spsc.Enqueue(i)
		if i < cap {
			assert.ErrorIs(t, err, nil)
		} else {
			assert.ErrorIs(t, err, ErrIsFull)
		}
	}

	assert.Equal(t, spsc.Length(), cap)
	assert.Equal(t, spsc.Capacity(), cap)
	for i := 0; i < cap*2; i++ {
		_, err := spsc.Dequeue()
		if i < cap {
			assert.ErrorIs(t, err, nil)
		} else {
			assert.ErrorIs(t, err, ErrIsEmpty)
		}
	}
	assert.Equal(t, spsc.Length(), 0)
	assert.Equal(t, spsc.Capacity(), cap)
}
