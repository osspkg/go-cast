/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"errors"
	"fmt"

	"go.osspkg.com/cast"
)

func main() {
	_, decodeErr := cast.StrToLimit[int]("12345", 4)
	_, encodeErr := cast.StringEncodeLimit("12345", 4)

	fmt.Println(errors.Is(decodeErr, cast.ErrSizeLimit))
	fmt.Println(errors.Is(encodeErr, cast.ErrSizeLimit))
}
