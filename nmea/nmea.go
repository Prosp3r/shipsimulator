package nmea

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/adrianmo/go-nmea"
)

type server struct {
	nmeaFile    string
	address     string
	delay       int
	loop        bool
	nmeaReadCh  chan string
	connections *connections
}

func NewServer(nmeaFile string, address string, delay int, loop bool) *server {
	s := server{
		nmeaFile:    nmeaFile,
		address:     address,
		delay:       delay,
		loop:        loop,
		nmeaReadCh:  make(chan string),
		connections: newConnections(),
	}

	return &s
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
func (s *server) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	nl, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Printf("error: net listen failed: %v\n", err)
		os.Exit(1)
	}
	defer nl.Close()

	// Wait for new connections, and put it on the newConn channel
	// to be added to the conn map.
	go func() {
		for {
			conn, err := nl.Accept()
			if err != nil {
				log.Printf("error: conn.Accept failed: %v\n", err)
			}
			go func() {
				s.connections.newConn <- conn
			}()
		}
	}()

	wg.Add(1)
	go func() {
		err := s.connections.handle(ctx)
		if err != nil {
			log.Printf("%v\n", err)
		}
		wg.Done()
	}()

	go s.readFile()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		// Block and wait for CTRL+C
		sig := <-sigCh
		fmt.Printf("Got exit signal, terminating all processes, %v\n", sig)

		cancel()

		wg.Wait()

		os.Exit(0)
	}()

	for v := range s.nmeaReadCh {
		s.connections.newData <- []byte(v)
	}
	// }

}

type connections struct {
	newConn chan net.Conn
	newData chan []byte
	conns   map[net.Conn]bool
}

func newConnections() *connections {
	c := connections{
		newConn: make(chan net.Conn),
		newData: make(chan []byte),
		conns:   make(map[net.Conn]bool),
	}
	return &c
}

func (c *connections) handle(ctx context.Context) error {
	for {
		select {
		case conn := <-c.newConn:
			c.conns[conn] = true

		case b := <-c.newData:
			var wg sync.WaitGroup

			// Channel to put the not-active connections to delete.
			deleteConn := make(chan net.Conn, len(c.conns))

			if len(c.conns) != 0 {

				for conn := range c.conns {
					wg.Add(1)
					go func(conn net.Conn) {
						defer wg.Done()

						// Check if connection is active, put it on the
						// delete channel if not active.
						tmpB := make([]byte, 1)
						conn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))

						_, err := conn.Read(tmpB)
						if err == io.EOF {
							log.Printf("error: connection timed out: %v\n", err)
							deleteConn <- conn
							return
						}

						// Connection active, write data to it.
						_, err = conn.Write(b)
						if err != nil {
							log.Printf("error: conn.Write: %v\n", err)
						}
					}(conn)
				}

				wg.Wait()

				// We are done checking what conn's that are active. Close
				// the channel so the range do not block.
				close(deleteConn)

				if len(deleteConn) != 0 {
					for v := range deleteConn {
						delete(c.conns, v)
					}
				}
			}

		case <-ctx.Done():
			return fmt.Errorf("info: connection.Handle: got done signal")
		}
	}
}

func (s *server) readFile() error {
	// Open the nmea file for reading
	f, err := os.Open(s.nmeaFile)
	if err != nil {
		return fmt.Errorf("error: failed to open nmea file for reading: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	ticker := time.NewTicker(time.Duration(s.delay) * time.Microsecond)

	for range ticker.C {
		// Check if there are more to scan
		for scanner.Scan() {
			sentence, err := nmea.Parse(scanner.Text())
			if err != nil {
				return fmt.Errorf("error: failed to parse nmea sentence: %v", err)
			}

			if sentence.DataType() == nmea.TypeRMC {
				rmc := sentence.(nmea.RMC).String()
				s.nmeaReadCh <- rmc
				break
			}

		}

	}
	return nil
}
