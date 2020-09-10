package main

import (
	"flag"
	"log"

	"github.com/RaaLabs/shipsimulator/nmea"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start().Stop()

	file := flag.String("file", "../output.nmea", "The name of the the NMEA file to read")
	address := flag.String("address", "localhost:8888", "The network host and port to send to, like localhost:8888")
	delay := flag.Int("delay", 1000, "The delay to wait between each send of data")
	flag.Parse()

	err := nmea.Run(*file, *address, *delay)
	if err != nil {
		log.Println(err)
	}
}
