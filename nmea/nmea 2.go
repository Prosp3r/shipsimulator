package nmea

import (
	"bufio"
	"fmt"
	"io"
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
		rmc := sentence.(nmea.RMC).String() + "\n"

		n, err := conn.Write([]byte(rmc))
		if err != nil && n != 0 {
			return fmt.Errorf("error: conn write failed sendToBroker: %v", err)
		}
	}
	return nil
}

// Run will start the parsing and sending process,
// and takes the full path of the file to parse,
// the address:port of the host to connect to,
// and a delay as an int in milliseconds to wait
// between each iteration of line in the file.
func Run(nmeaFile string, address string, delay int, loop bool) error {

	// Open the network connection to the receiver
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("error: net dial failed sendToBroker: %v", err)
	}
	defer conn.Close()

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
	}

	return nil
}
