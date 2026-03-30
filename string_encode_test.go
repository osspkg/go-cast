package cast_test

import (
	"bytes"
	"encoding/xml"
	"errors"
	"testing"
	"time"

	"go.osspkg.com/cast"
)

// Вспомогательные структуры для тестирования интерфейсов
type testByter []byte

func (b testByter) Bytes() []byte {
	return []byte(b)
}

type testStringer string

func (s testStringer) String() string {
	return string(s)
}

type testGoStringer string

func (s testGoStringer) GoString() string {
	return "Go" + string(s)
}

type testBinaryMarshaler string

func (b testBinaryMarshaler) MarshalBinary() ([]byte, error) {
	return []byte(b), nil
}

type testTextMarshaler string

func (t testTextMarshaler) MarshalText() ([]byte, error) {
	return []byte(t), nil
}

type testJSONMarshaler string

func (j testJSONMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(j) + `"`), nil
}

type testXMLStruct struct {
	Value string `xml:"value"`
}

func (x *testXMLStruct) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(x.Value, start)
}

type testError string

func (e testError) Error() string {
	return string(e)
}

type testUnsupported struct{}

// Тестовые структуры для проверки JSON и XML маршалинга
type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type testMap map[string]int

// Вспомогательные структуры для тестирования ошибок
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

type errorBinaryMarshaler struct {
	err error
}

func (b *errorBinaryMarshaler) MarshalBinary() ([]byte, error) {
	return nil, b.err
}

type errorTextMarshaler struct {
	err error
}

func (t *errorTextMarshaler) MarshalText() ([]byte, error) {
	return nil, t.err
}

func TestUnit_StringEncode(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	duration := 2 * time.Second

	tests := []struct {
		name        string
		input       any
		expected    string
		expectError bool
	}{
		// nil case
		{
			name:     "nil",
			input:    nil,
			expected: "",
		},

		// io.Reader case
		{
			name:     "io.Reader",
			input:    bytes.NewBufferString("reader content"),
			expected: "reader content",
		},

		// Byter interface
		{
			name:     "Byter",
			input:    testByter("byter content"),
			expected: "byter content",
		},

		// Stringer interface
		{
			name:     "Stringer",
			input:    testStringer("stringer content"),
			expected: "stringer content",
		},

		// GoStringer interface
		{
			name:     "GoStringer",
			input:    testGoStringer("gostringer content"),
			expected: "Gogostringer content",
		},

		// BinaryMarshaler interface
		{
			name:     "BinaryMarshaler",
			input:    testBinaryMarshaler("binary marshaler content"),
			expected: "binary marshaler content",
		},

		// TextMarshaler interface
		{
			name:     "TextMarshaler",
			input:    testTextMarshaler("text marshaler content"),
			expected: "text marshaler content",
		},

		// JSONMarshaler interface
		{
			name:     "JSONMarshaler",
			input:    testJSONMarshaler("json marshaler content"),
			expected: `"json marshaler content"`,
		},

		// XML case (using xml.Marshal for structs)
		{
			name:     "XML struct",
			input:    &testXMLStruct{Value: "xml content"},
			expected: `<testXMLStruct>xml content</testXMLStruct>`,
		},

		// Error interface
		{
			name:     "error",
			input:    testError("error content"),
			expected: "error content",
		},

		// Basic types
		{
			name:     "string",
			input:    "basic string",
			expected: "basic string",
		},
		{
			name:     "[]byte",
			input:    []byte("byte slice"),
			expected: "byte slice",
		},
		{
			name:     "int",
			input:    42,
			expected: "42",
		},
		{
			name:     "int8",
			input:    int8(42),
			expected: "42",
		},
		{
			name:     "int16",
			input:    int16(42),
			expected: "42",
		},
		{
			name:     "int32",
			input:    int32(42),
			expected: "42",
		},
		{
			name:     "int64",
			input:    int64(42),
			expected: "42",
		},
		{
			name:     "uint",
			input:    uint(42),
			expected: "42",
		},
		{
			name:     "uint8",
			input:    uint8(42),
			expected: "42",
		},
		{
			name:     "uint16",
			input:    uint16(42),
			expected: "42",
		},
		{
			name:     "uint32",
			input:    uint32(42),
			expected: "42",
		},
		{
			name:     "uint64",
			input:    uint64(42),
			expected: "42",
		},
		{
			name:     "float32",
			input:    float32(3.14),
			expected: "3.14",
		},
		{
			name:     "float64",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "bool true",
			input:    true,
			expected: "true",
		},
		{
			name:     "bool false",
			input:    false,
			expected: "false",
		},

		// Time types
		{
			name:     "time.Duration",
			input:    duration,
			expected: "2s",
		},
		{
			name:     "time.Time",
			input:    now,
			expected: now.Format(time.RFC3339),
		},

		// Complex types
		{
			name:     "struct",
			input:    testStruct{Name: "John", Age: 30},
			expected: `{"name":"John","age":30}`,
		},
		{
			name:     "map",
			input:    testMap{"a": 1, "b": 2},
			expected: `{"a":1,"b":2}`,
		},
		{
			name:     "slice",
			input:    []int{1, 2, 3},
			expected: `[1,2,3]`,
		},
		{
			name:     "array",
			input:    [3]int{1, 2, 3},
			expected: `[1,2,3]`,
		},

		// Pointers
		{
			name:     "pointer to string",
			input:    func() *string { s := "pointer string"; return &s }(),
			expected: "pointer string",
		},
		{
			name:     "pointer to struct",
			input:    &testStruct{Name: "Doe", Age: 25},
			expected: `{"name":"Doe","age":25}`,
		},

		// Unsupported types
		{
			name:        "unsupported func",
			input:       func() {},
			expectError: true,
		},
		{
			name:        "unsupported chan",
			input:       make(chan int),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cast.StringEncode(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("result mismatch\nexpected: %q\nactual:   %q", tt.expected, result)
			}
		})
	}

	// Тест для проверки обработки ошибок в io.Reader
	t.Run("io.Reader error", func(t *testing.T) {
		errReader := &errorReader{err: errors.New("read error")}
		_, err := cast.StringEncode(errReader)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if err.Error() != "read error" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// Тест для проверки обработки ошибок в BinaryMarshaler
	t.Run("BinaryMarshaler error", func(t *testing.T) {
		errMarshaler := &errorBinaryMarshaler{err: errors.New("binary marshal error")}
		_, err := cast.StringEncode(errMarshaler)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if err.Error() != "binary marshal error" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// Дополнительный тест для проверки, что nil-указатель обрабатывается как nil
	t.Run("nil pointer", func(t *testing.T) {
		var ptr *string
		result, err := cast.StringEncode(ptr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "" {
			t.Errorf("expected empty string for nil pointer, got %q", result)
		}
	})

	// Тест для проверки, что структура без поддерживаемых интерфейсов маршалится в JSON
	t.Run("struct without interfaces", func(t *testing.T) {
		type simpleStruct struct {
			Field1 string
			Field2 int
		}

		s := simpleStruct{
			Field1: "test",
			Field2: 42,
		}

		expected := `{"Field1":"test","Field2":42}`
		result, err := cast.StringEncode(s)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("result mismatch\nexpected: %q\nactual:   %q", expected, result)
		}
	})
}
