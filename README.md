# go-cast

[![Go Version](https://img.shields.io/github/go-mod/go-version/osspkg/go-cast)](https://go.dev/)
[![License](https://img.shields.io/github/license/osspkg/go-cast)](LICENSE)
[![CI](https://github.com/osspkg/go-cast/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/osspkg/go-cast/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/osspkg/go-cast)](https://goreportcard.com/report/github.com/osspkg/go-cast)
[![Go Reference](https://pkg.go.dev/badge/go.osspkg.com/cast.svg)](https://pkg.go.dev/go.osspkg.com/cast)

`go-cast` is a standard-library-only Go package for converting strings and Go
values. It supports generic conversions, common encoding interfaces, JSON/XML
fallbacks, and bounded input/output sizes.

## Demo

```go
package main

import (
	"fmt"
	"log"

	"go.osspkg.com/cast"
)

func main() {
	value, err := cast.StrTo[int]("42")
	if err != nil {
		log.Fatal(err)
	}

	encoded, err := cast.StringEncode(value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("value=%d encoded=%s\n", value, encoded)
}
```

Output:

```text
value=42 encoded=42
```

## Getting Started

The package requires Go 1.26 or newer:

```bash
go get go.osspkg.com/cast
```

Import the package using its module path:

```go
import "go.osspkg.com/cast"
```

## Features

### String to value

- `StrTo[T]` converts a string to a generic Go value.
- `StrToSlice[T]` converts a separated string to `[]T`.
- `StringDecode` decodes a string into a destination pointer.
- Supported destinations include strings, byte slices, numbers, booleans,
  `time.Duration`, `time.Time`, writers, and standard encoding interfaces.
- Structs, maps, arrays, and slices use JSON decoding when no more specific
  handler applies.

### Value to string

`StringEncode` supports strings, byte slices, numeric and boolean values,
durations, timestamps, pointers, readers, errors, custom `Stringer` and
`Byter` implementations, and standard marshaling interfaces. Structs, maps,
arrays, and slices use JSON encoding as a fallback.

### Size limits

Conversions use a default 16 MiB limit defined by `DefaultMaxBytes`. Use the
`*Limit` variants when a different non-negative limit is required:

```go
value, err := cast.StrToLimit[int]("42", 1024)
text, err := cast.StringEncodeLimit(value, 1024)
```

When a conversion exceeds the configured limit, the returned error matches
`cast.ErrSizeLimit` through `errors.Is`.

Reader conversions are bounded while reading. A reader that repeatedly returns
`(0, nil)` is rejected with `io.ErrNoProgress`.

For an empty input, `StringDecodeLimit` clears `string` and `[]byte`
destinations. Other destinations remain unchanged and are not initialized.
`StrToLimit` returns the target type's zero value without invoking
`Initializer`.

## Contributing

See [AGENTS.md](AGENTS.md) for the repository architecture, API invariants,
development workflow, and guidance for coding agents.

Run the standard checks before opening a pull request:

```bash
go test ./...
go test -race ./...
go vet ./...
```

The repository Makefile also provides `make ci`, which runs the project
workflow through `goppy`. The `license` and linting targets may update files;
review the working tree after running them.

## Contributors

See the [GitHub contributors graph](https://github.com/osspkg/go-cast/graphs/contributors).

## License

This project is distributed under the BSD 3-Clause License. See
[LICENSE](LICENSE) for the full text.
