package main

import (
	"flag"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	localAddress := flag.String("localAddress", "127.0.0.1:8080", "ipAddress:port")
	remoteAddress := flag.String("remoteAddress", "erter.org:80", "remoteAddress:port")
	flag.Parse()

	listener, err := net.Listen("tcp", *localAddress)
	if err != nil {
		log.Printf("error: net.Listen: %v\n", err)
		return
	}
	defer listener.Close()

	for {
		lc, err := listener.Accept()
		if err != nil {
			log.Printf("error: net.Listen: %v\n", err)
			return
		}

		go connectRemote(lc, *remoteAddress)
	}

}

// connectRemote will connect to the remote,
// and copy all the input/output between the
// connections.
// Will close the connection when done.
func connectRemote(lc net.Conn, remoteAddress string) {
	rc, err := net.Dial("tcp", remoteAddress)
	if err != nil {
		log.Printf("error: net.Dial: %v\n", err)
		return
	}

	// Send data from local to remote site
	go func() {
		defer rc.Close()
		defer lc.Close()
		_, err := io.Copy(rc, lc)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("error: io.Copy remote<-local : %v\n", err)
			return
		}
	}()

	// Receive data from remote back to local
	go func() {
		defer rc.Close()
		defer lc.Close()
		_, err := io.Copy(lc, rc)
		if err != nil {
			log.Printf("error: io.Copy local<-remote : %v\n", err)
			return
		}

	}()

}
