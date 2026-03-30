/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cast

type Initializer interface {
	Initialize() error
}

type Byter interface {
	Bytes() []byte
}

type Stringer interface {
	String() string
}

type UnStringer interface {
	UnString(string)
}
