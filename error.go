package ringbuffer

import "errors"

var (
	// ErrIsEmpty indicate ring buffer is empty
	ErrIsEmpty = errors.New("ring buffer is empty")
	// ErrIsFull indicate ring buffer is full
	ErrIsFull = errors.New("ring buffer is full")
	// ErrElementIsNil will return when the element are enqueue is nil
	ErrElementIsNil = errors.New("enqueued element should not be nil")
	// ErrSlotIsReading When element are enqueued, the elements in the same slot have not been consumed
	ErrSlotIsReading = errors.New("consumer is processing at current slot")
	// ErrSlotIsWriting When an element has not been fully enqueued, the consumer has started to consume
	ErrSlotIsWriting = errors.New("producer is processing at current slot")
)

// IsInProcessingError return true when error is indicate the element enqueue or dequeue not been fully complete
func IsInProcessingError(err error) bool {
	return err == ErrSlotIsReading || err == ErrSlotIsWriting
}
