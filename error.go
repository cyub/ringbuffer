// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ringbuffer

import "errors"

var (
	// ErrIsEmpty indicate ring buffer is empty
	ErrIsEmpty = errors.New("ring buffer is empty")
	// ErrIsFull indicate ring buffer is full
	ErrIsFull = errors.New("ring buffer is full")
	// ErrSlotIsReading When element are enqueued, the elements in the same slot have not been consumed
	ErrSlotIsReading = errors.New("consumer is processing at current slot")
	// ErrSlotIsWriting When an element has not been fully enqueued, the consumer has started to consume
	ErrSlotIsWriting = errors.New("producer is processing at current slot")
)

// IsInProcessingError return true when error is indicate the element enqueue or dequeue not been fully complete
func IsInProcessingError(err error) bool {
	return err == ErrSlotIsReading || err == ErrSlotIsWriting
}
