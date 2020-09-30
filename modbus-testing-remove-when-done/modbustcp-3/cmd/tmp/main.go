package main

import (
	"fmt"
	"math"
)

func main() {
	n := uint16((math.Float32bits(3.1415)) & 0xffff)

	fmt.Printf("%b\n", n)
}
