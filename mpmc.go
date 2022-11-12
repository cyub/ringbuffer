// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"sync/atomic"
	"unsafe"
)

// MpmcRingBuffer define Multi-Producer/Multi-Consumer ring buffer
type MpmcRingBuffer struct {
	head     uint64
	tail     uint64
	capacity int
	elements []interface{}
}

var _ RingBuffer = (*MpmcRingBuffer)(nil)

// NewMpmcRingBuffer return the mpmc ring buffer with specified capacity
func NewMpmcRingBuffer(capacity int) *MpmcRingBuffer {
	return &MpmcRingBuffer{
		head:     0,
		tail:     0,
		capacity: capacity,
		elements: make([]interface{}, capacity),
	}
}

// Enqueue element to the ring buffer
// if the ring buffer is full, then return ErrIsFull. if the enqueue elem is nil, return ErrElementIsNil
// When equeue element to the ring buffer, it may happen that the consumer who are consuming the same slot, so an error is returned
func (q *MpmcRingBuffer) Enqueue(elem interface{}) error {
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
		if atomic.LoadPointer(&slot.typ) != nil {
			return ErrSlotIsReading
		}
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
func (q *MpmcRingBuffer) Dequeue() (interface{}, error) {
	for {
		h := atomic.LoadUint64(&q.head)
		t := atomic.LoadUint64(&q.tail)
		if t == h {
			return nil, ErrIsEmpty
		}

		slot := (*eface)(unsafe.Pointer(&q.elements[h%uint64(q.capacity)]))
		if atomic.LoadPointer(&slot.val) == nil {
			return nil, ErrSlotIsWriting
		}
		if atomic.CompareAndSwapUint64(&q.head, h, h+1) {
			elem := *(*interface{})(unsafe.Pointer(slot))
			atomic.StorePointer(&slot.val, nil)
			atomic.StorePointer(&slot.typ, nil)
			return elem, nil
		}
	}
}

// Length return the number of all elements
func (q *MpmcRingBuffer) Length() int {
	return int(atomic.LoadUint64(&q.tail) - atomic.LoadUint64(&q.head))
}

// Capacity return the capacity of ring buffer
func (q *MpmcRingBuffer) Capacity() int {
	return q.capacity
}
