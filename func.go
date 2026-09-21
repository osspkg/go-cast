/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cast

// WriteFunc adapts a function to the io.Writer interface.
type WriteFunc func([]byte) (int, error)

// Write implements io.Writer.
func (w WriteFunc) Write(p []byte) (int, error) {
	return w(p)
}

// ReadFunc adapts a function to the io.Reader interface.
type ReadFunc func([]byte) (int, error)

// Read implements io.Reader.
func (r ReadFunc) Read(p []byte) (int, error) {
	return r(p)
}
