package main

import (
	"flag"
	"log"

	"github.com/RaaLabs/shipsimulator/nmea"
	// "github.com/pkg/profile"
)

func main() {
	// defer profile.Start().Stop()

	file := flag.String("file", "./output.nmea", "The name of the the NMEA file to read")
	address := flag.String("address", "localhost:8888", "The network host and port to send to, like localhost:8888")
	delay := flag.Int("delay", 1000000, "The delay to wait between each send of data given in Micro Seconds. Default is 1000000 (1 Second)")
	loop := flag.Bool("loop", false, "loop over again, and again, and again, and again,...........")
	flag.Parse()

	err := nmea.Run(*file, *address, *delay, *loop)
	if err != nil {
		log.Println(err)
	}
}
