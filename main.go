package main

import (
	"fmt"
)

func main() {
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)

	for i := 5; i >= 0; i-- {
		result, _ := Divide(100, i)
		fmt.Printf("i = %d, 100/%d = %d\n", i, i, result)
	}
}

func Divide(divisible, divider int) (int, error) {
	return divisible / divider, nil
}
