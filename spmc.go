// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"sync/atomic"
	"unsafe"
)

// SpmcRingBuffer define Single Producer/Multi-Consumer ring buffer
type SpmcRingBuffer struct {
	head     uint64
	tail     uint64
	capacity int
	elements []interface{}
}

var _ RingBuffer = (*SpmcRingBuffer)(nil)

// NewSpmcRingBuffer return the spmc ring buffer with specified capacity
func NewSpmcRingBuffer(capacity int) *SpmcRingBuffer {
	return &SpmcRingBuffer{
		head:     0,
		tail:     0,
		capacity: capacity,
		elements: make([]interface{}, capacity),
	}
}

// Enqueue element to the ring buffer
// if the ring buffer is full, then return ErrIsFull.
// When equeue element to the ring buffer, it may happen that the consumer who are consuming the same slot, so an error is returned
func (q *SpmcRingBuffer) Enqueue(elem interface{}) error {
	if elem == nil {
		elem = nilPlaceholder
	}
retry:
	h := atomic.LoadUint64(&q.head)
	t := q.tail
	if t == (h + uint64(q.capacity)) {
		return ErrIsFull
	}

	slot := (*eface)(unsafe.Pointer(&q.elements[t%uint64(q.capacity)]))
	if atomic.LoadPointer(&slot.typ) != nil {
		goto retry
	}

	*(*interface{})(unsafe.Pointer(slot)) = elem
	atomic.StoreUint64(&q.tail, t+1)
	return nil
}

// Dequeue an element from the ring buffer
// if the ring buffer is empty, then return ErrIsEmpty
func (q *SpmcRingBuffer) Dequeue() (interface{}, error) {
	for {
		h := atomic.LoadUint64(&q.head)
		t := atomic.LoadUint64(&q.tail)
		if t == h {
			return nil, ErrIsEmpty
		}
		if atomic.CompareAndSwapUint64(&q.head, h, h+1) {
			slot := (*eface)(unsafe.Pointer(&q.elements[h%uint64(q.capacity)]))
			elem := *(*interface{})(unsafe.Pointer(slot))
			slot.val = nil
			atomic.StorePointer(&slot.typ, nil)
			if elem == nilPlaceholder {
				return nil, nil
			}
			return elem, nil
		}
	}
}

// Length return the number of all elements
func (q *SpmcRingBuffer) Length() int {
	return int(atomic.LoadUint64(&q.tail) - atomic.LoadUint64(&q.head))
}

// Capacity return the capacity of ring buffer
func (q *SpmcRingBuffer) Capacity() int {
	return q.capacity
}
