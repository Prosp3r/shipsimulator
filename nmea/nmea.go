package nmea

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/adrianmo/go-nmea"
)

// readData will read the file by name provided,
// and call the nmea parser line by line.
func readData(conn net.Conn, f io.Reader) error {

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		err := parse(scanner.Text(), conn)
		if err != nil {
			return err
		}
	}

	return nil
}

func parse(nmeaText string, conn net.Conn) error {
	sentence, err := nmea.Parse(nmeaText)
	if err != nil {
		return fmt.Errorf("error: failed to parse nmea sentence: %v", err)
	}

	if sentence.DataType() == nmea.TypeRMC {
		rmc := sentence.(nmea.RMC).String() + "\n"

		n, err := conn.Write([]byte(rmc))
		if err != nil && n != 0 {
			return fmt.Errorf("error: conn write failed sendToBroker: %v", err)
		}
	}
	return nil
}

func Run(nmeaFile string, address string) error {

	// Open the network connection to the receiver
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("error: net dial failed sendToBroker: %v", err)
	}
	defer conn.Close()

	// Open the nmea file for reading
	f, err := os.Open(nmeaFile)
	if err != nil {
		return fmt.Errorf("error: failed to open nmea file for reading: %v", err)
	}
	defer f.Close()

	// Start the reading and sending
	err = readData(conn, f)
	if err != nil {
		return err
	}

	return nil
}
