// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"sync/atomic"
)

// SpscRingBuffer define Single Producer/Single Consumer ring buffer
type SpscRingBuffer struct {
	head     uint64
	tail     uint64
	capacity int
	elements []interface{}
}

var _ RingBuffer = (*SpscRingBuffer)(nil)

// NewSpscRingBuffer return the spsc ring buffer with specified capacity
func NewSpscRingBuffer(capacity int) *SpscRingBuffer {
	return &SpscRingBuffer{
		head:     0,
		tail:     0,
		capacity: capacity,
		elements: make([]interface{}, capacity),
	}
}

// Enqueue element to the ring buffer
// if the ring buffer is full, then return ErrIsFull
func (q *SpscRingBuffer) Enqueue(elem interface{}) error {
	if elem == nil {
		elem = nilPlaceholder
	}
	h := atomic.LoadUint64(&q.head)
	t := q.tail
	if t == (h + uint64(q.capacity)) {
		return ErrIsFull
	}

	q.elements[t%uint64(q.capacity)] = elem
	atomic.AddUint64(&q.tail, 1)
	return nil
}

// Dequeue an element from the ring buffer
// if the ring buffer is empty, then return ErrIsEmpty
func (q *SpscRingBuffer) Dequeue() (interface{}, error) {
	h := q.head
	t := atomic.LoadUint64(&q.tail)
	if t == h {
		return nil, ErrIsEmpty
	}

	elem := q.elements[h%uint64(q.capacity)]
	atomic.AddUint64(&q.head, 1)
	if elem == nilPlaceholder {
		return nil, nil
	}
	return elem, nil
}

// Length return the number of all elements
func (q *SpscRingBuffer) Length() int {
	return int(atomic.LoadUint64(&q.tail) - atomic.LoadUint64(&q.head))
}

// Capacity return the capacity of ring buffer
func (q *SpscRingBuffer) Capacity() int {
	return q.capacity
}
