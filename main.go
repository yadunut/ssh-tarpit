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
		for i, conn := range connections {
			if conn == nil {
				continue
			}
			if _, err := conn.Write([]byte{'a', '\n'}); err != nil {
				log.Printf("Closing Connection %d from: %s\n", count, conn.RemoteAddr().String())
				count--
				conn.Close()
				connections[i] = nil
			}
		}
		time.Sleep(DELAY)
	}
}

// func PrintMemUsage() {
// var m runtime.MemStats
// runtime.ReadMemStats(&m)
// // For info on each, see: https://golang.org/pkg/runtime/#MemStats

// log.Printf("\tTotalAlloc = %v KiB", bToKb(m.TotalAlloc))
// log.Printf("\tSys = %v KiB", bToKb(m.Sys))
// log.Printf("\tNumGC = %v\n", m.NumGC)
// }

// func bToKb(b uint64) uint64 {
// return b / 1024
// }
