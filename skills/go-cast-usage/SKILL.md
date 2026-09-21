---
name: go-cast-usage
description: Use the go-cast library for bounded string/value conversions in Go, choosing the correct API, limits, error handling, and custom interfaces.
---

# go-cast Usage

Use this skill when writing or reviewing Go code that imports `go.osspkg.com/cast`,
or when explaining how to convert strings and Go values with this library.

## Package boundary

- Import the package as `go.osspkg.com/cast`.
- The library is standard-library-only and targets Go 1.26.
- Prefer the public conversion functions and interfaces; do not recreate the
  conversion dispatch with reflection or a second parsing layer.

## Choose the API

- String to one value: `StrTo[T]` or `StrToLimit[T]`.
- Separated string to a slice: `StrToSlice[T]` or `StrToSliceLimit[T]`.
- Decode into an existing destination: `StringDecode` or `StringDecodeLimit`.
- Value to string: `StringEncode` or `StringEncodeLimit`.
- Use a `*Limit` function whenever the caller has a resource-specific size
  budget; otherwise the default is 16 MiB.

Read [references/api_reference.md](references/api_reference.md) for the complete
dispatch rules and [references/usage_patterns.md](references/usage_patterns.md)
for focused recipes. The runnable examples are under `examples/`.

## Invariants to preserve

- Limits must be non-negative. A conversion over the configured limit matches
  `ErrSizeLimit` with `errors.Is`.
- Reader conversion is bounded and returns `io.ErrNoProgress` for a reader
  that reports `(0, nil)`.
- Empty `StringDecode` input clears `string` and `[]byte` destinations. Other
  destinations remain unchanged and are not initialized.
- Empty `StrTo[T]` input returns `T`'s zero value without invoking
  `Initializer`.
- Primitive encoding uses the library's decimal and shortest `strconv` formats;
  do not change output formatting casually.

## Implementation guidance

1. Choose the narrowest public function that matches the conversion direction.
2. Pass an explicit limit when input can be supplied by a user, network, file,
   or untrusted integration.
3. Handle and return errors; use `errors.Is(err, cast.ErrSizeLimit)` for limit
   failures.
4. Add a focused regression test for new conversion behavior, especially
   malformed input, limits, short writes, reader errors, and empty input.
5. Keep examples compilable and import the real module path.
