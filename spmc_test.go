// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"runtime"
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

func BenchmarkRingSPMC(b *testing.B) {
	b.Run("1P100C", func(b *testing.B) {
		benchmarkRingBuffer(b, NewSpmcRingBuffer(8192), 100, false, true, doneEnv())
	})
	b.Run("1P4C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkRingBuffer(b, NewSpmcRingBuffer(8192), 4, false, true, doneEnv())
		runtime.GOMAXPROCS(pp)
	})
}

func BenchmarkChanSPMC(b *testing.B) {
	b.Run("1P100C", func(b *testing.B) {
		benchmarkRingBuffer(b, newChanRingBuffer(8192), 100, false, true, doneEnv())
	})
	b.Run("1P4C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkRingBuffer(b, newChanRingBuffer(8192), 4, false, true, doneEnv())
		runtime.GOMAXPROCS(pp)
	})
}
