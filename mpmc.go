// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// MpmcRingBuffer define Multi-Producer/Multi-Consumer ring buffer
type MpmcRingBuffer struct {
	ringbuffer
}

var _ RingBuffer = (*MpmcRingBuffer)(nil)

// NewMpmcRingBuffer return the mpmc ring buffer with specified capacity
func NewMpmcRingBuffer(capacity int) *MpmcRingBuffer {
	return &MpmcRingBuffer{
		ringbuffer{
			head:     0,
			tail:     0,
			capacity: capacity,
			elements: make([]interface{}, capacity),
		},
	}
}

// Enqueue element to the ring buffer
// if the ring buffer is full, then return ErrIsFull.
// When equeue element to the ring buffer, it may happen that the consumer who are consuming the same slot, so an error is returned
func (q *MpmcRingBuffer) Enqueue(elem interface{}) error {
	if elem == nil {
		elem = nilPlaceholder
	}

	for {
		h := atomic.LoadUint64(&q.head)
		t := atomic.LoadUint64(&q.tail)
		if t >= h+uint64(q.capacity) {
			return ErrIsFull
		}

		slot := (*eface)(unsafe.Pointer(&q.elements[t%uint64(q.capacity)]))
		if atomic.LoadPointer(&slot.typ) != nil {
			runtime.Gosched()
			continue
		}
		if !atomic.CompareAndSwapUint64(&q.tail, t, t+1) {
			runtime.Gosched()
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
			runtime.Gosched()
			continue
		}
		if !atomic.CompareAndSwapUint64(&q.head, h, h+1) {
			runtime.Gosched()
			continue
		}
		elem := *(*interface{})(unsafe.Pointer(slot))
		atomic.StorePointer(&slot.val, nil)
		atomic.StorePointer(&slot.typ, nil)
		if elem == nilPlaceholder {
			return nil, nil
		}
		return elem, nil
	}
}
