package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	b, err := ioutil.ReadFile("./sample.bin")
	if err != nil {
		log.Printf("error : %v\n", err)
	}

	fmt.Printf("%v", b)
}
