# go-cast

> Standard-library-only Go conversions between strings and values.

## Overview

`go-cast` is the `cast` package published as `go.osspkg.com/cast`. It converts
strings to typed values and values to strings through explicit type handling,
standard encoding interfaces, and JSON/XML fallbacks.

The module targets Go 1.26 and currently has no third-party Go dependencies.
The package does not start background goroutines or own long-lived resources.

## Source of truth

- `go.mod` defines the module path and supported Go version.
- `README.md` is the user-facing overview and usage guide.
- Exported GoDoc comments are the API reference rendered by `pkg.go.dev`.
- `string_decode.go` contains string-to-value conversion behavior.
- `string_encode.go` contains value-to-string conversion behavior.
- `limits.go` contains shared size-limit definitions and helpers.
- `AGENTS.md` documents repository constraints for coding agents.

## Quick start

```bash
go get go.osspkg.com/cast
```

```go
package main

import (
	"fmt"

	"go.osspkg.com/cast"
)

func main() {
	value, err := cast.StrTo[int]("42")
	if err != nil {
		panic(err)
	}
	text, err := cast.StringEncode(value)
	if err != nil {
		panic(err)
	}
	fmt.Println(value, text)
}
```

## API reference

### Conversion functions

- `StrTo[T]` converts a string to `T` using `DefaultMaxBytes`.
- `StrToLimit[T]` converts a string to `T` with an explicit byte limit.
- `StrToSlice[T]` converts a separated string to `[]T`.
- `StrToSliceLimit[T]` adds an explicit byte limit to slice conversion.
- `StringDecode` decodes into a destination pointer using the default limit.
- `StringDecodeLimit` decodes into a destination pointer with an explicit limit.
- `StringEncode` converts a value to a string using the default limit.
- `StringEncodeLimit` converts a value to a string with an explicit limit.

### Public types and errors

- `DefaultMaxBytes` is the default 16 MiB input/output limit.
- `ErrSizeLimit` identifies a conversion that exceeded its byte limit; callers
  should use `errors.Is`.
- `Initializer` is called before non-empty string decoding.
- `Byter`, `Stringer`, and `UnStringer` provide custom conversion hooks.
- `ReadFunc` and `WriteFunc` adapt functions to reader/writer interfaces.
- `Ptr` returns a pointer to a value.

## Behavioral invariants

- Explicit limits must be non-negative. A conversion exceeding its limit
  returns an error matching `ErrSizeLimit`.
- Default conversions are bounded by `DefaultMaxBytes`.
- `StringDecodeLimit` validates the destination before handling empty input.
- Empty input clears `string` and `[]byte` destinations. Other destinations
  remain unchanged, and `Initializer` is not called.
- `StrToLimit` returns the target type's zero value for empty input without
  invoking `Initializer`.
- `StringEncode` handles nil and typed nil pointers as an empty string.
- Reader conversion is bounded and returns `io.ErrNoProgress` when a reader
  reports no progress.
- Primitive numeric and boolean encoding uses `strconv`; preserve the existing
  textual formats when changing conversion code.

## Development workflow

Run focused checks directly with the Go toolchain:

```bash
go test ./...
go test -race ./...
go vet ./...
```

The Makefile provides the repository workflow:

```bash
make install
make ci
```

`make install` bootstraps `goppy`. `make ci` runs license generation, linting,
tests, and the build through the project tooling. Inspect `git diff` after
running it because license and linting steps may rewrite files.

## Change guidance

- Keep the module path `go.osspkg.com/cast` and the Go 1.26 target.
- Prefer the standard library; do not add a dependency for functionality
  already provided by `strconv`, `encoding`, `io`, `reflect`, or `time`.
- Preserve the default and explicit size-limit contracts.
- Add regression tests for behavior changes, malformed input, short writes,
  reader errors, and resource-limit boundaries.
- Update `README.md` and exported GoDoc when public behavior changes.
- Keep changes focused and do not reformat unrelated files.

## Validation expectations

Before handing off a change, run formatting, unit tests, race tests, vet, the
configured linter, and `git diff --check`. Report exactly which checks were
run and any environment-dependent checks that could not be completed.
