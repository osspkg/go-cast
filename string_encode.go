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

// StringEncode converts obj to a string using DefaultMaxBytes as the output limit.
func StringEncode(obj any) (string, error) {
	return StringEncodeLimit(obj, DefaultMaxBytes)
}

// StringEncodeLimit converts obj to a string using maxBytes as the output limit.
func StringEncodeLimit(obj any, maxBytes int) (s string, err error) {
	if err = checkSize(0, maxBytes); err != nil {
		return "", err
	}
	return stringEncode(obj, maxBytes)
}

func stringEncode(obj any, maxBytes int) (s string, err error) {
	if obj == nil {
		return
	}

	ref := reflect.ValueOf(obj)
	if ref.Kind() == reflect.Pointer && ref.IsNil() {
		return
	}

	switch v := obj.(type) {

	case string:
		return limitString(v, maxBytes)

	case []byte:
		return bytesToString(v, maxBytes)

	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool:
		return limitString(fmt.Sprintf("%v", v), maxBytes)

	case time.Duration:
		return limitString(v.String(), maxBytes)

	case time.Time:
		return limitString(v.Format(time.RFC3339), maxBytes)

	case io.Reader:
		b, readErr := readAllLimit(v, maxBytes)
		if readErr != nil {
			return "", readErr
		}
		return bytesToString(b, maxBytes)

	case Byter:
		return bytesToString(v.Bytes(), maxBytes)

	case Stringer:
		return limitString(v.String(), maxBytes)

	case fmt.GoStringer:
		return limitString(v.GoString(), maxBytes)

	case encoding.BinaryMarshaler:
		b, marshalErr := v.MarshalBinary()
		if marshalErr != nil {
			return "", marshalErr
		}
		return bytesToString(b, maxBytes)

	case encoding.TextMarshaler:
		b, marshalErr := v.MarshalText()
		if marshalErr != nil {
			return "", marshalErr
		}
		return bytesToString(b, maxBytes)

	case json.Marshaler:
		b, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			return "", marshalErr
		}
		return bytesToString(b, maxBytes)

	case xml.Marshaler:
		b, marshalErr := xml.Marshal(v)
		if marshalErr != nil {
			return "", marshalErr
		}
		return bytesToString(b, maxBytes)

	case error:
		return limitString(v.Error(), maxBytes)

	default:
		switch ref.Kind() {
		case reflect.Pointer:
			return stringEncode(ref.Elem().Interface(), maxBytes)

		case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
			b, marshalErr := json.Marshal(obj)
			if marshalErr != nil {
				return "", marshalErr
			}
			return bytesToString(b, maxBytes)

		default:
			err = fmt.Errorf("unsupported type: %T", obj)
		}
	}

	return
}

func readAllLimit(r io.Reader, maxBytes int) ([]byte, error) {
	const chunkSize = 32 << 10

	capacity := chunkSize
	if maxBytes < capacity {
		capacity = maxBytes
	}
	data := make([]byte, 0, capacity)
	chunk := make([]byte, chunkSize)

	for {
		n, err := r.Read(chunk)
		if n < 0 || n > len(chunk) {
			return nil, fmt.Errorf("cast: reader returned invalid byte count %d", n)
		}
		if sizeErr := checkSize(len(data)+n, maxBytes); sizeErr != nil {
			return nil, sizeErr
		}
		data = append(data, chunk[:n]...)
		if err == io.EOF {
			return data, nil
		}
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, io.ErrNoProgress
		}
	}
}
