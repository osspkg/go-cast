package main

import (
	"errors"
	"fmt"

	"go.osspkg.com/cast"
)

func main() {
	_, decodeErr := cast.StrToLimit[int]("12345", 4)
	_, encodeErr := cast.StringEncodeLimit("12345", 4)

	fmt.Println(errors.Is(decodeErr, cast.ErrSizeLimit))
	fmt.Println(errors.Is(encodeErr, cast.ErrSizeLimit))
}
