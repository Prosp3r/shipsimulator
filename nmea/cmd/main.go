package main

import (
	"flag"
	"log"

	"github.com/RaaLabs/shipsimulator/nmea"
)

func main() {
	file := flag.String("file", "../output.nmea", "The name of the the NMEA file to read")
	address := flag.String("address", "localhost:8888", "The network host and port to send to, like localhost:8888")
	flag.Parse()

	err := nmea.Run(*file, *address)
	if err != nil {
		log.Println(err)
	}
}
