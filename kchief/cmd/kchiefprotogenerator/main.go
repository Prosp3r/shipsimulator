package main

import (
	"log"

	"github.com/RaaLabs/shipsimulator/kchief"
)

func main() {
	err := kchief.RunProtoGenerator()
	if err != nil {
		log.Printf("%v\n", err)
	}

}
