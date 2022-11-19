// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSPSCRingBuffer(t *testing.T) {
	cap := 16
	spsc := NewSpscRingBuffer(cap)

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

func BenchmarkRingSPSC(b *testing.B) {
	b.Run("1P1C", func(b *testing.B) {
		benchmarkRingBuffer(b, NewSpscRingBuffer(8192), 100, false, false, doneEnv())
	})
	b.Run("1P1C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkRingBuffer(b, NewSpscRingBuffer(8192), 4, false, false, doneEnv())
		runtime.GOMAXPROCS(pp)
	})
}

func BenchmarkChanSPSC(b *testing.B) {
	b.Run("1P1C", func(b *testing.B) {
		benchmarkRingBuffer(b, newChanRingBuffer(8192), 100, false, false, doneEnv())
	})
	b.Run("1P1C_1CPU", func(b *testing.B) {
		pp := runtime.GOMAXPROCS(1)
		benchmarkRingBuffer(b, newChanRingBuffer(8192), 4, false, false, doneEnv())
		runtime.GOMAXPROCS(pp)
	})
}
