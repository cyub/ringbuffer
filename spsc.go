package ringbuffer

import (
	"sync/atomic"
)

type spscRingBuffer struct {
	head     uint64
	tail     uint64
	capacity int
	elements []interface{}
}

var _ RingBuffer = (*spscRingBuffer)(nil)

func NewSpscRingBuffer(capacity int) *spscRingBuffer {
	return &spscRingBuffer{
		head:     0,
		tail:     0,
		capacity: capacity,
		elements: make([]interface{}, capacity),
	}
}

func (q *spscRingBuffer) Enqueue(val interface{}) error {
	h := atomic.LoadUint64(&q.head)
	t := q.tail
	if t == (h + uint64(q.capacity)) {
		return ErrIsFull
	}

	q.elements[t%uint64(q.capacity)] = val
	atomic.AddUint64(&q.tail, 1)
	return nil
}

func (q *spscRingBuffer) Dequeue() (interface{}, error) {
	h := q.head
	t := atomic.LoadUint64(&q.tail)
	if t == h {
		return nil, ErrIsEmpty
	}

	val := q.elements[h%uint64(q.capacity)]
	atomic.AddUint64(&q.head, 1)
	return val, nil
}

func (q *spscRingBuffer) Length() int {
	return int(atomic.LoadUint64(&q.tail) - atomic.LoadUint64(&q.head))
}

func (q *spscRingBuffer) Capacity() int {
	return q.capacity
}
