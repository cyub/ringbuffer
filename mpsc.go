// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"sync/atomic"
	"unsafe"
)

// MpscRingBuffer define Multi-Producer/Single Consumer ring buffer
type MpscRingBuffer struct {
	head     uint64
	tail     uint64
	capacity int
	elements []interface{}
}

var _ RingBuffer = (*MpscRingBuffer)(nil)

// NewMpscRingBuffer return the mpsc ring buffer with specified capacity
func NewMpscRingBuffer(capacity int) *MpscRingBuffer {
	return &MpscRingBuffer{
		head:     0,
		tail:     0,
		capacity: capacity,
		elements: make([]interface{}, capacity),
	}
}

// Enqueue element to the ring buffer
// if the ring buffer is full, then return ErrIsFull. if the enqueue elem is nil, return ErrElementIsNil
func (q *MpscRingBuffer) Enqueue(elem interface{}) error {
	if elem == nil {
		return ErrElementIsNil
	}

	for {
		h := atomic.LoadUint64(&q.head)
		t := atomic.LoadUint64(&q.tail)
		if t == h+uint64(q.capacity) {
			return ErrIsFull
		}

		slot := (*eface)(unsafe.Pointer(&q.elements[t%uint64(q.capacity)]))
		if !atomic.CompareAndSwapUint64(&q.tail, t, t+1) {
			continue
		}

		eval := *(*eface)(unsafe.Pointer(&elem))
		atomic.StorePointer(&slot.typ, eval.typ)
		atomic.StorePointer(&slot.val, eval.val)
		return nil
	}
}

// Dequeue an element from the ring buffer
// if the ring buffer is empty, then return ErrIsEmpty
// When dequeue element to the ring buffer, it may happen that the producer who are working the same slot, so an error is returned
func (q *MpscRingBuffer) Dequeue() (interface{}, error) {
	h := q.head
	t := atomic.LoadUint64(&q.tail)
	if t == h {
		return nil, ErrIsEmpty
	}

	slot := (*eface)(unsafe.Pointer(&q.elements[h%uint64(q.capacity)]))
	if atomic.LoadPointer(&slot.val) == nil {
		return nil, ErrSlotIsWriting
	}
	elem := *(*interface{})(unsafe.Pointer(slot))
	atomic.AddUint64(&q.head, 1)
	return elem, nil
}

// Length return the number of all elements
func (q *MpscRingBuffer) Length() int {
	return int(atomic.LoadUint64(&q.tail) - atomic.LoadUint64(&q.head))
}

// Capacity return the capacity of ring buffer
func (q *MpscRingBuffer) Capacity() int {
	return q.capacity
}
