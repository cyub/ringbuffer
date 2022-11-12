// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import (
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
