// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"sync/atomic"
	"unsafe"
)

// RingBuffer is an interface
type RingBuffer interface {
	Enqueue(interface{}) error
	Dequeue() (interface{}, error)
	Length() int
	Capacity() int
}

// eface empty interface
type eface struct {
	typ unsafe.Pointer
	val unsafe.Pointer
}

type placeholder struct{}

var nilPlaceholder interface{} = placeholder{}

type cacheLinePad struct {
	_ [128 - unsafe.Sizeof(uint64(0))%128]byte
}

// ringbuffer struct
type ringbuffer struct {
	head     uint64
	_        cacheLinePad
	tail     uint64
	_        cacheLinePad
	capacity int
	elements []interface{}
}

// Length return the number of all elements
func (q *ringbuffer) Length() int {
	return int(atomic.LoadUint64(&q.tail) - atomic.LoadUint64(&q.head))
}

// Capacity return the capacity of ring buffer
func (q *ringbuffer) Capacity() int {
	return q.capacity
}
