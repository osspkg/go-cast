/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cast

import (
	"errors"
	"fmt"
)

// DefaultMaxBytes is the default input or output limit used by conversions.
const DefaultMaxBytes = 16 << 20

// ErrSizeLimit reports that a conversion exceeded its configured byte limit.
var ErrSizeLimit = errors.New("cast: size limit exceeded")

func checkSize(size, maxBytes int) error {
	if maxBytes < 0 {
		return fmt.Errorf("cast: maximum bytes must be non-negative")
	}
	if size > maxBytes {
		return fmt.Errorf("%w: %d bytes exceeds %d-byte limit", ErrSizeLimit, size, maxBytes)
	}
	return nil
}

func bytesToString(data []byte, maxBytes int) (string, error) {
	if err := checkSize(len(data), maxBytes); err != nil {
		return "", err
	}
	return string(data), nil
}

func limitString(value string, maxBytes int) (string, error) {
	if err := checkSize(len(value), maxBytes); err != nil {
		return "", err
	}
	return value, nil
}
