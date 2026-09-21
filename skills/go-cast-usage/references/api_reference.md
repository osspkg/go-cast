# go-cast API reference

This reference describes the public API of the module `go.osspkg.com/cast`.
Use the package's GoDoc as the final authority if the implementation changes.

## Installation

```bash
go get go.osspkg.com/cast
```

```go
import "go.osspkg.com/cast"
```

## Conversion functions

| Function | Direction | Use |
| --- | --- | --- |
| `StrTo[T](s)` | string -> `T` | Default 16 MiB input limit. |
| `StrToLimit[T](s, maxBytes)` | string -> `T` | Explicit input limit. |
| `StrToSlice[T](s, sep)` | separated string -> `[]T` | Default 16 MiB input limit. |
| `StrToSliceLimit[T](s, sep, maxBytes)` | separated string -> `[]T` | Explicit input limit. |
| `StringDecode(obj, s)` | string -> destination | Decodes into a non-nil pointer with the default limit. |
| `StringDecodeLimit(obj, s, maxBytes)` | string -> destination | Decodes into a pointer with an explicit limit. |
| `StringEncode(obj)` | value -> string | Default 16 MiB output limit. |
| `StringEncodeLimit(obj, maxBytes)` | value -> string | Explicit output limit. |

## Decoding dispatch

`StringDecode` and `StringDecodeLimit` accept a non-nil pointer. They support:

- `*string` and `*[]byte`;
- all signed and unsigned integer types;
- `*float32`, `*float64`, `*complex64`, and `*complex128`;
- `*bool`, `*time.Duration`, and `*time.Time`;
- `io.Writer` and `io.StringWriter`;
- `UnStringer`;
- `encoding.BinaryUnmarshaler`, `encoding.TextUnmarshaler`,
  `json.Unmarshaler`, and `xml.Unmarshaler`;
- struct, map, array, and slice pointers through `encoding/json` as a fallback.

For non-empty input, a destination implementing `Initializer` is initialized
before dispatch. An empty input is handled after pointer validation and before
initialization: string and byte-slice destinations are cleared, while other
destinations are left unchanged.

Parsing uses the standard library. Numeric and boolean parse errors are
returned to the caller; malformed JSON/XML and custom unmarshaler errors are
also returned.

## Encoding dispatch

`StringEncode` and `StringEncodeLimit` support:

- strings and byte slices;
- signed and unsigned integers, floating-point values, and booleans;
- `time.Duration` and `time.Time` formatted as RFC 3339;
- `io.Reader`, `Byter`, `Stringer`, `fmt.GoStringer`, and `error`;
- `encoding.BinaryMarshaler`, `encoding.TextMarshaler`,
  `json.Marshaler`, and XML marshaling;
- pointers, structs, maps, arrays, and slices through JSON fallback.

Dispatch follows the type-switch order in the implementation. If a value
implements several supported interfaces, the earlier matching case wins.
Nil and typed nil pointers encode as an empty string.

Primitive formatting is decimal for integers and booleans use `true`/`false`.
Floating-point values use `strconv.FormatFloat` with format `g`, precision `-1`,
and the value's bit size.

## Limits and errors

`DefaultMaxBytes` is 16 MiB. Every default conversion is bounded by it.
Explicit `maxBytes` values must be non-negative. When the input or output
exceeds the selected limit, the error matches `cast.ErrSizeLimit`:

```go
if errors.Is(err, cast.ErrSizeLimit) {
	// Reject or report oversized input/output.
}
```

Reader conversion reads in bounded chunks. A reader that returns `(0, nil)` is
rejected with `io.ErrNoProgress`; a reader must eventually return data, an
error, or `io.EOF`.

## Related examples

- [`examples/basic/main.go`](../examples/basic/main.go) — generic conversion
  and JSON fallback encoding.
- [`examples/limits/main.go`](../examples/limits/main.go) — explicit limits and
  `errors.Is`.
- [`examples/custom/main.go`](../examples/custom/main.go) — `UnStringer` and
  `Byter` hooks.
