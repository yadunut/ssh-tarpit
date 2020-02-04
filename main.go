package main

import (
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// ADDRESS to run server on
const (
	ADDRESS = ":2222"
	DELAY   = time.Second * 3
)

var (
	connections = make([]Conn, 0, 100)
	count       = 0
)

// Conn defines a connection
type Conn struct {
	conn  net.Conn
	begin time.Time
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

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

		log.WithFields(log.Fields{
			"IP":         conn.RemoteAddr().String(),
			"count":      count,
			"connection": "received",
		}).Println()
		connections = append(connections, Conn{conn, time.Now()})
	}
}

func work() {
	for {
		for i := 0; i < len(connections); i++ {
			conn := connections[i]
			if conn.conn == nil {
				connections = append(connections[:i], connections[i+1:]...)
				continue
			}

			if _, err := conn.conn.Write([]byte{'a'}); err != nil {

				log.WithFields(log.Fields{
					"IP":         conn.conn.RemoteAddr().String(),
					"count":      count,
					"connection": "closed",
					"time-in":    time.Since(conn.begin).Seconds(),
				}).Println()
				count--
				conn.conn.Close()
				connections = append(connections[:i], connections[i+1:]...)
			}
		}
		time.Sleep(DELAY)

	}
	// Using range to iterate through the array cannot be used cause connections is modified during the loop execution. Range copies the connections array and hits a NPE
}
