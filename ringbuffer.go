package ringbuffer

import (
	"errors"
)

var (
	ErrIsEmpty = errors.New("ring buffer is empty")
	ErrIsFull  = errors.New("ring buffer is full")
)

type RingBuffer interface {
	Enqueue(interface{}) error
	Dequeue() (interface{}, error)
	Length() int
	Capacity() int
}
