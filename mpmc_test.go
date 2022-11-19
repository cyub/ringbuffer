// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"os"
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

func doneEnv() bool {
	switch os.Getenv("DONE") {
	case "", "true":
		return true
	case "false":
		return false
	default:
		return true
	}
}

func benchmarkRingBuffer(b *testing.B, ring RingBuffer, workers int, multiProducer, multiConsumer bool, done bool) {
	var size = b.N/workers + 1
	var numbers = mknumslice(size)
	var wg sync.WaitGroup

	wg.Add(workers * 2)
	b.ResetTimer()

	if multiProducer {
		for p := 0; p < workers; p++ {
			go func() {
				for i := range numbers {
				retry:
					if err := ring.Enqueue(i); err != nil {
						if done { // all must done
							runtime.Gosched()
							goto retry
						}
					}
				}
				wg.Done()
			}()
		}
	} else {
		go func() {
			for p := 0; p < workers; p++ {
				for i := range numbers {
				retry:
					if err := ring.Enqueue(i); err != nil {
						if done {
							runtime.Gosched()
							goto retry
						}
					}
				}
				wg.Done()
			}
		}()

	}

	if multiConsumer {
		for p := 0; p < workers; p++ {
			go func() {
				for i := 0; i < size; i++ {
				retry:
					if _, err := ring.Dequeue(); err != nil {
						if done {
							runtime.Gosched()
							goto retry
						}

					}
				}
				wg.Done()
			}()
		}
	} else {
		go func() {
			for p := 0; p < workers; p++ {
				for i := 0; i < size; i++ {
				retry:
					if _, err := ring.Dequeue(); err != nil {
						if done {
							runtime.Gosched()
							goto retry
						}
					}
				}
				wg.Done()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkRingMPMC(b *testing.B) {
	b.Run("100P100C", func(b *testing.B) {
		benchmarkRingBuffer(b, NewMpmcRingBuffer(8192), 100, true, true, doneEnv())
	})
	b.Run("4P4C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkRingBuffer(b, NewMpmcRingBuffer(8192), 4, true, true, doneEnv())
		runtime.GOMAXPROCS(pp)
	})
}

func BenchmarkChanMPMC(b *testing.B) {
	b.Run("100P100C", func(b *testing.B) {
		benchmarkRingBuffer(b, newChanRingBuffer(8192), 100, true, true, doneEnv())
	})
	b.Run("4P4C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkRingBuffer(b, newChanRingBuffer(8192), 4, true, true, doneEnv())
		runtime.GOMAXPROCS(pp)
	})
}

// ringbuffer base channel as compare object
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
