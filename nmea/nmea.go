package nmea

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/adrianmo/go-nmea"
)

// readData will read the file by name provided,
// and call the parse function for each line read.
// A time ticker have been used to schedule the reading of
// lines from the file at given intervals.
func readData(conn net.Conn, f io.Reader, delay int) error {

	scanner := bufio.NewScanner(f)
	ticker := time.NewTicker(time.Duration(delay) * time.Microsecond)

	for range ticker.C {
		// Check if there are more to scan
		if scanner.Scan() {
			err := parse(scanner.Text(), conn)
			if err != nil {
				return err
			}
		} else {
			break
		}

	}

	return nil

}

// parse will parse the NMEA data out of the text line read.
func parse(nmeaText string, conn net.Conn) error {
	sentence, err := nmea.Parse(nmeaText)
	if err != nil {
		return fmt.Errorf("error: failed to parse nmea sentence: %v", err)
	}

	// Send over the network connection to the receiver
	// if the data read is of the correct type.
	if sentence.DataType() == nmea.TypeRMC {
		rmc := sentence.(nmea.RMC).String()

		n, err := conn.Write([]byte(rmc))
		if err != nil && n != 0 {
			return fmt.Errorf("error: conn write failed sendToBroker: %v", err)
		}
	}
	return nil
}

// Run will start the parsing and sending process.
// Takes the "full path" of the file to parse.
// The "address:port" of the host to connect to in
// "mode=send", and the "address:port" of a local
// interface to listen on if "mode=listen".
// "delay" as an int in milliseconds to wait
// between each iteration of line in the file.
// Loop set to true will read the input file over
// and over.
func Run(nmeaFile string, address string, delay int, loop bool) {
	nl, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("error: net listen failed: %v\n", err)
		return
	}
	defer nl.Close()

	for {
		conn, err := nl.Accept()
		if err != nil {
			log.Printf("error: net Accept failed: %v\n", err)
		}

		go readAndSend(nmeaFile, conn, delay, loop)
	}

}
func readAndSend(nmeaFile string, conn net.Conn, delay int, loop bool) error {
	for {
		// Open the nmea file for reading
		f, err := os.Open(nmeaFile)
		if err != nil {
			return fmt.Errorf("error: failed to open nmea file for reading: %v", err)
		}

		// Start the reading and sending
		err = readData(conn, f, delay)
		if err != nil {
			return err
		}

		f.Close()

		if !loop {
			break
		}

		err = conn.Close()
		if err != nil {
			return fmt.Errorf("error: failed to close connection: %v", err)
		}

	}

	return nil
}
