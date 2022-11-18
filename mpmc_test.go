// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMPMCRingBuffer(t *testing.T) {
	cap := 16
	spsc := NewMpmcRingBuffer(cap)

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

func mknumslice(n int) []int {
	var s = make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func benchmarkMPMC(b *testing.B, ring RingBuffer, workers int, done bool) {
	var size = b.N/workers + 1
	var numbers = mknumslice(size)
	var wg sync.WaitGroup

	wg.Add(workers * 2)
	b.ResetTimer()
	for p := 0; p < workers; p++ {
		go func(p int) {
			for i := range numbers {
			try:
				err := ring.Enqueue(i)
				if err != nil {
					if done {
						runtime.Gosched()
						goto try
					}
				}
			}
			wg.Done()
		}(p)
	}

	for p := 0; p < workers; p++ {
		go func(c int) {
			for i := 0; i < size; i++ {
			try:
				_, err := ring.Dequeue()
				if err != nil {
					if done {
						runtime.Gosched()
						goto try
					}
				}

			}
			wg.Done()
		}(p)
	}
	wg.Wait()
}

func BenchmarkRingMPMC(b *testing.B) {
	b.Run("100P100C", func(b *testing.B) {
		benchmarkMPMC(b, NewMpmcRingBuffer(8192), 100, false)
	})
	b.Run("4P4C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkMPMC(b, NewMpmcRingBuffer(8192), 4, false)
		runtime.GOMAXPROCS(pp)
	})

	b.Run("100P100C_DONE", func(b *testing.B) {
		benchmarkMPMC(b, NewMpmcRingBuffer(8192), 100, true)
	})
	b.Run("4P4C_1CPU_DONE", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkMPMC(b, NewMpmcRingBuffer(8192), 4, true)
		runtime.GOMAXPROCS(pp)
	})
}

func BenchmarkChanMPMC(b *testing.B) {
	b.Run("100P100C", func(b *testing.B) {
		benchmarkMPMC(b, newChanRingBuffer(8192), 100, false)
	})
	b.Run("4P4C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkMPMC(b, newChanRingBuffer(8192), 4, false)
		runtime.GOMAXPROCS(pp)
	})

	b.Run("100P100C_DONE", func(b *testing.B) {
		benchmarkMPMC(b, newChanRingBuffer(8192), 100, true)
	})
	b.Run("4P4C_1CPU_DONE", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkMPMC(b, newChanRingBuffer(8192), 4, true)
		runtime.GOMAXPROCS(pp)
	})
}

// ringbuffer base channel
type chanRingBuffer chan interface{}

func newChanRingBuffer(capacity int) chanRingBuffer {
	return make(chanRingBuffer, capacity)
}
func (c chanRingBuffer) Enqueue(elem interface{}) error {
	select {
	case c <- elem:
		return nil
	default:
		return ErrIsFull
	}
}

func (c chanRingBuffer) Dequeue() (interface{}, error) {
	select {
	case elem := <-c:
		return elem, nil
	default:
		return nil, ErrIsEmpty
	}
}

func (c chanRingBuffer) Length() int {
	return len(c)
}

func (c chanRingBuffer) Capacity() int {
	return cap(c)
}
