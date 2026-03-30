/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cast

type WriteFunc func([]byte) (int, error)

func (w WriteFunc) Write(p []byte) (int, error) {
	return w(p)
}

type ReadFunc func([]byte) (int, error)

func (r ReadFunc) Read(p []byte) (int, error) {
	return r(p)
}
