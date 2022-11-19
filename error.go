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
)
