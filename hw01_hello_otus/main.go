package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	test := "Hello, OTUS!"
	fmt.Printf("%s\n", stringutil.Reverse(test))
}
