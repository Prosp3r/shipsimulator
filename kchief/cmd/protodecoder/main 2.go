/*
	Decode a file containing protobuf binary data into
	it's protobuf structure
*/

package main

import (
	"flag"
	"log"

	"github.com/RaaLabs/shipsimulator/kchief"
)

func main() {
	fileName := flag.String("fileName", "./sample/sample2.bin", "The full path with the filename of the Kchief protobuf data file")
	flag.Parse()

	if err := kchief.RunProtoReader(*fileName); err != nil {
		log.Println(err)
	}
}
