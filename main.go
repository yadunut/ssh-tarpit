package main

import (
	"log"
	"net"
	"os"
	"time"
)

const (
	ADDRESS = ":2222"
	DELAY   = time.Second * 3
)

var (
	connections = make([]net.Conn, 100)
)

func main() {
	log.Println("Starting")

	// Create Socket
	ll, err := net.Listen("tcp", ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	defer ll.Close()

	// Handles clean exits
	go func() {
		bb := make([]byte, 1)
		os.Stdin.Read(bb)
		os.Exit(1)
	}()

	go work()

	for {
		conn, err := ll.Accept()
		if err != nil {
			continue
		}
		log.Println("Received Connection from: ", conn.LocalAddr().String())
		connections = append(connections, conn)

	}
}

func work() {
	for {
		for i, conn := range connections {
			if conn == nil {
				continue
			}
			if _, err := conn.Write([]byte{'a', '\n'}); err != nil {
				log.Println("Closing Connection")
				conn.Close()
				connections[i] = nil
			}
		}
		time.Sleep(DELAY)
	}
}
