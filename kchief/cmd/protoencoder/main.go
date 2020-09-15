package main

import (
	"log"

	"github.com/RaaLabs/shipsimulator/kchief"
)

// TODO:
// - Specify output file
// - Read tags to create for from JSON file to replace the hard coded values in the protogenerator.go file

func main() {
	err := kchief.RunProtoEncode("myfile.bin")
	if err != nil {
		log.Printf("%v\n", err)
	}

}
