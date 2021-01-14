package main

import (
	"errors"
	"fmt"
	"reflect"
)

func main() {
	e1 := errors.New("Inner Error MSG")

	e2 := fmt.Errorf("> %w", e1)
	e3 := fmt.Errorf("> %w", e2)
	e4 := fmt.Errorf("> %w", e3)

	fmt.Println("....", e1)
	fmt.Println("....", e4)
	fmt.Println("....", errors.Unwrap(e4))
	fmt.Println("....", reflect.TypeOf(e4))
}
