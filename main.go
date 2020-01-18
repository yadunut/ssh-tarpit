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
	connections = make([]net.Conn, 0, 100)
	count       = 0
)

func main() {
	log.Println("Starting")

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
		count++

		log.Printf("Received Connection %d from: %s \n", count, conn.RemoteAddr().String())
		connections = append(connections, conn)

	}
}

func work() {
	for {
		for i := 0; i < len(connections); i++ {
			conn := connections[i]
			if conn == nil {
				continue
			}

			if _, err := conn.Write([]byte{'a', '\n'}); err != nil {
				log.Printf("Closing Connection %d from: %s\n", count, conn.RemoteAddr().String())
				count--
				conn.Close()
				connections = append(connections[:i], connections[i+1:]...)
			}
		}
		time.Sleep(DELAY)

	}
	// Using range to iterate through the array cannot be used cause connections is modified during the loop execution. Range copies the connections array and hits a NPE
}
