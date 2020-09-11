package main

import (
	"log"

	"github.com/RaaLabs/shipsimulator/kchief"
)

func main() {
	if err := kchief.Run(); err != nil {
		log.Println(err)
	}
}
