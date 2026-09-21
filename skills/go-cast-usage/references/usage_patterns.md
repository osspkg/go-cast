# go-cast usage patterns

Read this reference when generating application code or reviewing a conversion
for correctness and resource use.

## Convert trusted small values

Use the generic helpers when the input is already bounded by the surrounding
protocol:

```go
port, err := cast.StrTo[int](rawPort)
if err != nil {
	return fmt.Errorf("parse port: %w", err)
}
```

For separated values, pass the separator explicitly:

```go
ids, err := cast.StrToSlice[int](rawIDs, ",")
if err != nil {
	return fmt.Errorf("parse ids: %w", err)
}
```

## Bound untrusted input

Use an explicit limit at the boundary where a request, file, or message enters
the application:

```go
value, err := cast.StrToLimit[int](raw, 4<<10)
if err != nil {
	if errors.Is(err, cast.ErrSizeLimit) {
		return fmt.Errorf("input is too large: %w", err)
	}
	return err
}
```

The limit is measured in bytes. Do not treat it as a character or element
count.

## Decode into an existing value

`StringDecode` requires a non-nil pointer:

```go
var enabled bool
if err := cast.StringDecode(&enabled, raw); err != nil {
	return err
}
```

For empty input, string and byte-slice destinations are cleared. Other
destinations retain their previous value, and `Initializer` is not called.
Document or test this behavior when empty input is meaningful in the caller's
protocol.

## Encode values and readers

```go
text, err := cast.StringEncode(value)
if err != nil {
	return fmt.Errorf("encode value: %w", err)
}
```

When encoding an `io.Reader`, select `StringEncodeLimit` so a large or
untrusted stream cannot grow without bound. A reader that repeatedly returns
`(0, nil)` is rejected rather than looped forever.

## Use custom interfaces only for domain behavior

Implement `UnStringer` when a type needs custom string decoding and `Byter` or
`Stringer` when it needs custom encoding. Keep the methods deterministic and
return errors from validation at the caller boundary when the interface allows
it.

```go
type Identifier struct {
	Value string
}

func (id *Identifier) UnString(s string) {
	id.Value = strings.TrimSpace(s)
}
```

See `examples/custom/main.go` for a complete runnable version.
