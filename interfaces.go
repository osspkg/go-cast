/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cast

// Initializer initializes a value before it is decoded.
type Initializer interface {
	Initialize() error
}

// Byter returns the byte representation of a value.
type Byter interface {
	Bytes() []byte
}

// Stringer returns the string representation of a value.
type Stringer interface {
	String() string
}

// UnStringer updates a value from its string representation.
type UnStringer interface {
	UnString(string)
}
