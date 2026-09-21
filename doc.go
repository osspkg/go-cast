/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package cast converts strings and Go values with bounded input and output.
//
// Conversions use DefaultMaxBytes unless an explicit limit is supplied through
// one of the package's *Limit functions.
package cast
