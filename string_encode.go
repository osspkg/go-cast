/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cast

import (
	"encoding"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"time"
)

func StringEncode(obj any) (s string, err error) {
	if obj == nil {
		return
	}

	ref := reflect.ValueOf(obj)
	if ref.Kind() == reflect.Ptr && ref.IsNil() {
		return
	}

	switch v := obj.(type) {

	case string:
		s = v

	case []byte:
		s = string(v)

	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool:
		s = fmt.Sprintf("%v", v)

	case time.Duration:
		s = v.String()

	case time.Time:
		s = v.Format(time.RFC3339)

	case io.Reader:
		var b []byte
		b, err = io.ReadAll(v)
		s = string(b)

	case Byter:
		s = string(v.Bytes())

	case Stringer:
		s = v.String()

	case fmt.GoStringer:
		s = v.GoString()

	case encoding.BinaryMarshaler:
		var b []byte
		b, err = v.MarshalBinary()
		s = string(b)

	case encoding.TextMarshaler:
		var b []byte
		b, err = v.MarshalText()
		s = string(b)

	case json.Marshaler:
		var b []byte
		b, err = v.MarshalJSON()
		s = string(b)

	case xml.Marshaler:
		var b []byte
		b, err = xml.Marshal(v)
		s = string(b)

	case error:
		s = v.Error()

	default:
		switch ref.Kind() {
		case reflect.Ptr:
			return StringEncode(ref.Elem().Interface())

		case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
			var b []byte
			b, err = json.Marshal(obj)
			s = string(b)

		default:
			err = fmt.Errorf("unsupported type: %T", obj)
		}
	}

	return
}
