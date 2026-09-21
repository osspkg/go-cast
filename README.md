# go-cast

Go types converter.

Conversions use a default 16 MiB input/output limit (`DefaultMaxBytes`). Use
`StrToLimit`, `StrToSliceLimit`, `StringDecodeLimit`, or `StringEncodeLimit`
when a different non-negative limit is required. `ErrSizeLimit` identifies a
limit violation.
