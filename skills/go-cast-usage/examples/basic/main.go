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

	values, err := cast.StrToSlice[int]("1,2,3", ",")
	if err != nil {
		log.Fatal(err)
	}

	encoded, err := cast.StringEncode(values)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value)
	fmt.Println(values)
	fmt.Println(encoded)
}
