// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMPSCRingBuffer(t *testing.T) {
	cap := 16
	spsc := NewMpscRingBuffer(cap)

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

func BenchmarkRingMPSC(b *testing.B) {
	b.Run("100P1C", func(b *testing.B) {
		benchmarkRingBuffer(b, NewMpscRingBuffer(8192), 100, true, false, doneEnv())
	})
	b.Run("4P1C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkRingBuffer(b, NewMpscRingBuffer(8192), 4, true, false, doneEnv())
		runtime.GOMAXPROCS(pp)
	})
}

func BenchmarkChanMPSC(b *testing.B) {
	b.Run("100P1C", func(b *testing.B) {
		benchmarkRingBuffer(b, newChanRingBuffer(8192), 100, true, false, doneEnv())
	})
	b.Run("4P1C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkRingBuffer(b, newChanRingBuffer(8192), 4, true, false, doneEnv())
		runtime.GOMAXPROCS(pp)
	})
}
