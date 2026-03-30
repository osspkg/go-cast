package cast_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.osspkg.com/cast"
)

type testInitializer struct {
	init bool
}

func (t *testInitializer) Initialize() error {
	t.init = true
	return nil
}

type testUnStringer struct {
	Value string
}

func (u *testUnStringer) UnString(s string) {
	u.Value = s
}

type testBinaryUnmarshaler struct {
	Data []byte
}

func (b *testBinaryUnmarshaler) UnmarshalBinary(data []byte) error {
	b.Data = append([]byte{}, data...)
	return nil
}

type testTextUnmarshaler struct {
	Text string
}

func (t *testTextUnmarshaler) UnmarshalText(text []byte) error {
	t.Text = string(text)
	return nil
}

type testJSONUnmarshaler struct {
	Value string
}

func (j *testJSONUnmarshaler) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &j.Value)
}

type testXMLUnmarshaler struct {
	Object struct {
		Value string `xml:"value"`
	}
}

func (x *testXMLUnmarshaler) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return d.DecodeElement(&x.Object, &start)
}

type testWriter struct {
	Buffer *bytes.Buffer
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	return w.Buffer.Write(p)
}

type testStringWriter struct {
	Buffer *bytes.Buffer
}

func (w *testStringWriter) WriteString(s string) (n int, err error) {
	return w.Buffer.WriteString(s)
}

type testComplexStruct struct {
	Pointer *string
	Slice   []int
	Map     map[string]string
}

func TestStringDecode(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	rfc3339Time := now.Format(time.RFC3339)

	tests := []struct {
		name        string
		obj         any
		input       string
		expected    any
		expectError bool
	}{
		// Тесты для базовых типов
		{
			name:     "string",
			obj:      new(string),
			input:    "test string",
			expected: "test string",
		},
		{
			name:     "int",
			obj:      new(int),
			input:    "42",
			expected: 42,
		},
		{
			name:     "int8",
			obj:      new(int8),
			input:    "127",
			expected: int8(127),
		},
		{
			name:     "int16",
			obj:      new(int16),
			input:    "32767",
			expected: int16(32767),
		},
		{
			name:     "int32",
			obj:      new(int32),
			input:    "2147483647",
			expected: int32(2147483647),
		},
		{
			name:     "int64",
			obj:      new(int64),
			input:    "9223372036854775807",
			expected: int64(9223372036854775807),
		},
		{
			name:     "uint",
			obj:      new(uint),
			input:    "4294967295",
			expected: uint(4294967295),
		},
		{
			name:     "uint8",
			obj:      new(uint8),
			input:    "255",
			expected: uint8(255),
		},
		{
			name:     "uint16",
			obj:      new(uint16),
			input:    "65535",
			expected: uint16(65535),
		},
		{
			name:     "uint32",
			obj:      new(uint32),
			input:    "4294967295",
			expected: uint32(4294967295),
		},
		{
			name:     "uint64",
			obj:      new(uint64),
			input:    "18446744073709551615",
			expected: uint64(18446744073709551615),
		},
		{
			name:     "float32",
			obj:      new(float32),
			input:    "3.14",
			expected: float32(3.14),
		},
		{
			name:     "float64",
			obj:      new(float64),
			input:    "3.1415926535",
			expected: 3.1415926535,
		},
		{
			name:     "bool true",
			obj:      new(bool),
			input:    "true",
			expected: true,
		},
		{
			name:     "bool false",
			obj:      new(bool),
			input:    "false",
			expected: false,
		},
		{
			name:     "duration",
			obj:      new(time.Duration),
			input:    "2s",
			expected: 2 * time.Second,
		},
		{
			name:     "time",
			obj:      new(time.Time),
			input:    rfc3339Time,
			expected: now,
		},
		{
			name:     "[]byte",
			obj:      new([]byte),
			input:    "bytes content",
			expected: []byte("bytes content"),
		},

		// Тесты для интерфейсов
		{
			name:  "io.Writer",
			obj:   &testWriter{Buffer: &bytes.Buffer{}},
			input: "writer content",
			expected: func() testWriter {
				w := &testWriter{Buffer: &bytes.Buffer{}}
				_, _ = w.Write([]byte("writer content"))
				return *w
			}(),
		},
		{
			name:  "io.StringWriter",
			obj:   &testStringWriter{Buffer: &bytes.Buffer{}},
			input: "string writer content",
			expected: func() testStringWriter {
				w := &testStringWriter{Buffer: &bytes.Buffer{}}
				_, _ = w.WriteString("string writer content")
				return *w
			}(),
		},
		{
			name:  "UnStringer",
			obj:   &testUnStringer{},
			input: "unstringer content",
			expected: testUnStringer{
				Value: "unstringer content",
			},
		},
		{
			name:  "BinaryUnmarshaler",
			obj:   &testBinaryUnmarshaler{},
			input: "binary unmarshal content",
			expected: testBinaryUnmarshaler{
				Data: []byte("binary unmarshal content"),
			},
		},
		{
			name:  "TextUnmarshaler",
			obj:   &testTextUnmarshaler{},
			input: "text unmarshal content",
			expected: testTextUnmarshaler{
				Text: "text unmarshal content",
			},
		},
		{
			name:  "JSONUnmarshaler",
			obj:   &testJSONUnmarshaler{},
			input: `"json unmarshal content"`,
			expected: testJSONUnmarshaler{
				Value: "json unmarshal content",
			},
		},
		{
			name:  "XMLUnmarshaler",
			obj:   &testXMLUnmarshaler{},
			input: `<testXMLUnmarshaler><value>xml unmarshal content</value></testXMLUnmarshaler>`,
			expected: testXMLUnmarshaler{
				Object: testXMLStruct{
					Value: "xml unmarshal content",
				},
			},
		},

		// Тесты для сложных типов через JSON
		{
			name:  "struct",
			obj:   &testStruct{},
			input: `{"name":"John","age":30}`,
			expected: testStruct{
				Name: "John",
				Age:  30,
			},
		},
		{
			name:  "map",
			obj:   &testMap{},
			input: `{"a":1,"b":2}`,
			expected: testMap{
				"a": 1,
				"b": 2,
			},
		},
		{
			name:  "slice",
			obj:   &[]int{},
			input: `[1,2,3]`,
			expected: []int{
				1, 2, 3,
			},
		},
		{
			name:  "complex struct",
			obj:   &testComplexStruct{},
			input: `{"Pointer":"test","Slice":[1,2,3],"Map":{"a":"b"}}`,
			expected: testComplexStruct{
				Pointer: func() *string { s := "test"; return &s }(),
				Slice:   []int{1, 2, 3},
				Map:     map[string]string{"a": "b"},
			},
		},

		// Тесты для инициализации
		{
			name:  "initializer",
			obj:   &testInitializer{},
			input: "{}",
			expected: testInitializer{
				init: true,
			},
		},

		// Негативные тесты
		{
			name:        "not a pointer",
			obj:         "not a pointer",
			input:       "test",
			expectError: true,
		},
		{
			name:        "nil pointer",
			obj:         (*string)(nil),
			input:       "test",
			expectError: true,
		},
		{
			name:        "invalid int",
			obj:         new(int),
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "invalid bool",
			obj:         new(bool),
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "invalid duration",
			obj:         new(time.Duration),
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "invalid time",
			obj:         new(time.Time),
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "unsupported type",
			obj:         new(func()),
			input:       "test",
			expectError: true,
		},
		{
			name:        "unsupported JSON",
			obj:         new(testStruct),
			input:       `{"name":"John","age":"invalid"}`,
			expectError: true,
		},
		{
			name:     "empty string",
			obj:      new(string),
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cast.StringDecode(tt.obj, tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Сравнение результатов
			result := reflect.ValueOf(tt.obj).Elem().Interface()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("result mismatch\nexpected: %#v\nactual:   %#v", tt.expected, result)
			}
		})
	}

	// Тест для проверки обработки ошибок в интерфейсах
	t.Run("error in interface", func(t *testing.T) {
		errUnmarshaler := &errorTextUnmarshaler{err: errors.New("unmarshal error")}
		err := cast.StringDecode(errUnmarshaler, "test")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if err.Error() != "unmarshal error" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// Тест для проверки обработки ошибок в JSON unmarshal
	t.Run("error in json unmarshal", func(t *testing.T) {
		var s testStruct
		err := cast.StringDecode(&s, `{"name":"John","age":"invalid"}`)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}

// Вспомогательные структуры для тестирования ошибок
type errorBinaryUnmarshaler struct {
	err error
}

func (b *errorBinaryUnmarshaler) UnmarshalBinary(data []byte) error {
	return b.err
}

type errorTextUnmarshaler struct {
	err error
}

func (t *errorTextUnmarshaler) UnmarshalText(text []byte) error {
	return t.err
}

type errorJSONUnmarshaler struct {
	err error
}

func (j *errorJSONUnmarshaler) UnmarshalJSON(data []byte) error {
	return j.err
}

type errorXMLUnmarshaler struct {
	err error
}

func (x *errorXMLUnmarshaler) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return x.err
}

func TestStrTo(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    any
		expectError bool
	}{
		// Базовые типы
		{
			name:     "string",
			input:    "test",
			expected: "test",
		},
		{
			name:     "int",
			input:    "42",
			expected: 42,
		},
		{
			name:     "int8",
			input:    "127",
			expected: int8(127),
		},
		{
			name:     "bool",
			input:    "true",
			expected: true,
		},
		{
			name:     "duration",
			input:    "2s",
			expected: 2 * time.Second,
		},
		{
			name:     "time",
			input:    time.Now().UTC().Truncate(time.Second).Format(time.RFC3339),
			expected: time.Now().UTC().Truncate(time.Second),
		},

		// Пустая строка
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Используем рефлексию для создания обобщенного вызова
			result := reflect.New(reflect.TypeOf(tt.expected)).Elem().Interface()

			switch result.(type) {
			case string:
				val, err := cast.StrTo[string](tt.input)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if val != tt.expected {
					t.Errorf("result mismatch\nexpected: %v\nactual: %v", tt.expected, val)
				}

			case int:
				val, err := cast.StrTo[int](tt.input)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if val != tt.expected {
					t.Errorf("result mismatch\nexpected: %v\nactual: %v", tt.expected, val)
				}

			case int8:
				val, err := cast.StrTo[int8](tt.input)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if val != tt.expected {
					t.Errorf("result mismatch\nexpected: %v\nactual: %v", tt.expected, val)
				}

			case bool:
				val, err := cast.StrTo[bool](tt.input)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if val != tt.expected {
					t.Errorf("result mismatch\nexpected: %v\nactual: %v", tt.expected, val)
				}

			case time.Duration:
				val, err := cast.StrTo[time.Duration](tt.input)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if val != tt.expected {
					t.Errorf("result mismatch\nexpected: %v\nactual: %v", tt.expected, val)
				}

			case time.Time:
				val, err := cast.StrTo[time.Time](tt.input)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				// Сравниваем без наносекунд, так как time.Parse может их потерять
				if !val.Truncate(time.Second).Equal(tt.expected.(time.Time).Truncate(time.Second)) {
					t.Errorf("result mismatch\nexpected: %v\nactual: %v", tt.expected, val)
				}

			default:
				t.Fatalf("unsupported test type: %T", tt.expected)
			}
		})
	}
}

func TestStrToSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		sep         string
		expected    any
		expectError bool
	}{
		// Тесты для строк
		{
			name:     "string slice",
			input:    "a,b,c",
			sep:      ",",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty string slice",
			input:    "",
			sep:      ",",
			expected: []string(nil),
		},

		// Тесты для чисел
		{
			name:     "int slice",
			input:    "1,2,3",
			sep:      ",",
			expected: []int{1, 2, 3},
		},
		{
			name:     "float slice",
			input:    "1.1,2.2,3.3",
			sep:      ",",
			expected: []float64{1.1, 2.2, 3.3},
		},

		// Тесты для bool
		{
			name:     "bool slice",
			input:    "true,false,true",
			sep:      ",",
			expected: []bool{true, false, true},
		},

		// Тесты для времени
		{
			name:  "duration slice",
			input: fmt.Sprintf("%s,%s,%s", "1s", "2s", "3s"),
			sep:   ",",
			expected: []time.Duration{
				1 * time.Second,
				2 * time.Second,
				3 * time.Second,
			},
		},

		// Негативные тесты
		{
			name:        "invalid int slice",
			input:       "1,invalid,3",
			sep:         ",",
			expectError: true,
			expected:    []int{},
		},
		{
			name:     "custom separator",
			input:    "a|b|c",
			sep:      "|",
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Используем рефлексию для определения типа слайса
			sliceType := reflect.TypeOf(tt.expected)
			elemType := sliceType.Elem()

			// Выполняем преобразование
			switch elemType.Kind() {
			case reflect.String:
				val, err := cast.StrToSlice[string](tt.input, tt.sep)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(val, tt.expected) {
					t.Errorf("result mismatch\nexpected: %#v\nactual:   %#v", tt.expected, val)
				}

			case reflect.Int:
				val, err := cast.StrToSlice[int](tt.input, tt.sep)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(val, tt.expected) {
					t.Errorf("result mismatch\nexpected: %#v\nactual:   %#v", tt.expected, val)
				}

			case reflect.Int64:
				val, err := cast.StrToSlice[time.Duration](tt.input, tt.sep)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(val, tt.expected) {
					t.Errorf("result mismatch\nexpected: %#v\nactual:   %#v", tt.expected, val)
				}

			case reflect.Float64:
				val, err := cast.StrToSlice[float64](tt.input, tt.sep)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(val, tt.expected) {
					t.Errorf("result mismatch\nexpected: %#v\nactual:   %#v", tt.expected, val)
				}

			case reflect.Bool:
				val, err := cast.StrToSlice[bool](tt.input, tt.sep)
				if tt.expectError {
					if err == nil {
						t.Fatal("expected error but got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(val, tt.expected) {
					t.Errorf("result mismatch\nexpected: %#v\nactual:   %#v", tt.expected, val)
				}

			case reflect.Struct:
				if elemType == reflect.TypeOf(time.Duration(0)) {
					val, err := cast.StrToSlice[time.Duration](tt.input, tt.sep)
					if tt.expectError {
						if err == nil {
							t.Fatal("expected error but got nil")
						}
						return
					}
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					if !reflect.DeepEqual(val, tt.expected) {
						t.Errorf("result mismatch\nexpected: %#v\nactual:   %#v", tt.expected, val)
					}
				}
			default:
				t.Fatalf("unsupported test type: %T", tt.expected)
			}
		})
	}
}
