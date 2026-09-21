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
	"strconv"
	"strings"
	"time"
)

// StrTo converts a string to T using DefaultMaxBytes as the input limit.
func StrTo[T any](s string) (T, error) {
	return StrToLimit[T](s, DefaultMaxBytes)
}

// StrToLimit converts a string to T using maxBytes as the input limit.
// For an empty string, it returns T's zero value without invoking Initializer.
func StrToLimit[T any](s string, maxBytes int) (T, error) {
	var v T
	if err := checkSize(len(s), maxBytes); err != nil {
		return v, err
	}
	if len(s) == 0 {
		return v, nil
	}
	err := StringDecodeLimit(&v, s, maxBytes)
	return v, err
}

// StrToSlice converts a separated string to []T using DefaultMaxBytes as the input limit.
func StrToSlice[T any](s, sep string) ([]T, error) {
	return StrToSliceLimit[T](s, sep, DefaultMaxBytes)
}

// StrToSliceLimit converts a separated string to []T using maxBytes as the input limit.
func StrToSliceLimit[T any](s, sep string, maxBytes int) ([]T, error) {
	if err := checkSize(len(s), maxBytes); err != nil {
		return nil, err
	}
	if s == "" {
		return nil, nil
	}

	count := strings.Count(s, sep) + 1
	result := make([]T, 0, count)

	for _, ss := range strings.Split(s, sep) {
		val, err := StrToLimit[T](ss, maxBytes)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}

	return result, nil
}

// StringDecode decodes s into obj using DefaultMaxBytes as the input limit.
func StringDecode(obj any, s string) error {
	return StringDecodeLimit(obj, s, DefaultMaxBytes)
}

// StringDecodeLimit decodes s into obj using maxBytes as the input limit.
// An empty input clears string and byte-slice destinations; other destinations
// remain unchanged and are not initialized.
func StringDecodeLimit(obj any, s string, maxBytes int) error {
	if err := checkSize(len(s), maxBytes); err != nil {
		return err
	}

	ref := reflect.ValueOf(obj)
	if ref.Kind() != reflect.Pointer {
		return fmt.Errorf("got not a pointer")
	}

	if ref.IsNil() {
		return fmt.Errorf("got nil pointer")
	}

	if len(s) == 0 {
		switch p := obj.(type) {
		case *string:
			*p = s
		case *[]byte:
			*p = []byte(s)
		}
		return nil
	}

	if in, ok := obj.(Initializer); ok {
		if err := in.Initialize(); err != nil {
			return err
		}
	}

	var err error
	switch p := obj.(type) {

	case *string:
		*p = s

	case *[]byte:
		*p = []byte(s)

	case *int:
		var val int64
		val, err = strconv.ParseInt(s, 10, strconv.IntSize)
		*p = int(val)

	case *int8:
		var val int64
		val, err = strconv.ParseInt(s, 10, 8)
		*p = int8(val)

	case *int16:
		var val int64
		val, err = strconv.ParseInt(s, 10, 16)
		*p = int16(val)

	case *int32:
		var val int64
		val, err = strconv.ParseInt(s, 10, 32)
		*p = int32(val)

	case *int64:
		*p, err = strconv.ParseInt(s, 10, 64)

	case *uint:
		var val uint64
		val, err = strconv.ParseUint(s, 10, strconv.IntSize)
		*p = uint(val)

	case *uint8:
		var val uint64
		val, err = strconv.ParseUint(s, 10, 8)
		*p = uint8(val)

	case *uint16:
		var val uint64
		val, err = strconv.ParseUint(s, 10, 16)
		*p = uint16(val)

	case *uint32:
		var val uint64
		val, err = strconv.ParseUint(s, 10, 32)
		*p = uint32(val)

	case *uint64:
		*p, err = strconv.ParseUint(s, 10, 64)

	case *float32:
		var val float64
		val, err = strconv.ParseFloat(s, 32)
		*p = float32(val)

	case *float64:
		*p, err = strconv.ParseFloat(s, 64)

	case *complex64:
		var val complex128
		val, err = strconv.ParseComplex(s, 64)
		*p = complex64(val)

	case *complex128:
		*p, err = strconv.ParseComplex(s, 128)

	case *bool:
		*p, err = strconv.ParseBool(s)

	case *time.Duration:
		*p, err = time.ParseDuration(s)

	case *time.Time:
		*p, err = time.Parse(time.RFC3339, s)

	case io.Writer:
		data := []byte(s)
		var n int
		n, err = p.Write(data)
		if err == nil && n != len(data) {
			err = io.ErrShortWrite
		}

	case io.StringWriter:
		var n int
		n, err = p.WriteString(s)
		if err == nil && n != len(s) {
			err = io.ErrShortWrite
		}

	case UnStringer:
		p.UnString(s)

	case encoding.BinaryUnmarshaler:
		err = p.UnmarshalBinary([]byte(s))

	case encoding.TextUnmarshaler:
		err = p.UnmarshalText([]byte(s))

	case json.Unmarshaler:
		err = p.UnmarshalJSON([]byte(s))

	case xml.Unmarshaler:
		err = xml.Unmarshal([]byte(s), p)

	default:

		switch ref.Elem().Kind() {
		case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
			err = json.Unmarshal([]byte(s), obj)

		default:
			err = fmt.Errorf("unsupported type: %T", obj)
		}
	}

	return err
}
