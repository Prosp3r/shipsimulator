package main

import (
	"flag"
	"log"

	"github.com/RaaLabs/shipsimulator/kchief"
)

// TODO:
// - Specify output file
// - Read tags to create for from JSON file to replace the hard coded values in the protogenerator.go file

func main() {
	outBinFile := flag.String("outFile", "myfile.bin", "specify the filename of the output file")
	inJsonFile := flag.String("inJsonFile", "", "specify the full path of the json file to take as input")
	flag.Parse()

	err := kchief.RunProtoEncode(*outBinFile, *inJsonFile)
	if err != nil {
		log.Printf("%v\n", err)
	}

}
