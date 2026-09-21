/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"log"

	"go.osspkg.com/cast"
)

type identifier struct {
	Value string
}

func (id *identifier) UnString(value string) {
	id.Value = value
}

type token string

func (t token) Bytes() []byte {
	return []byte(t)
}

func main() {
	var id identifier
	if err := cast.StringDecode(&id, "user-42"); err != nil {
		log.Fatal(err)
	}

	encoded, err := cast.StringEncode(token(id.Value))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(id.Value)
	fmt.Println(encoded)
}
