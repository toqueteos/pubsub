package main

import (
	"fmt"
	"net"
	"os"

	"github.com/ikanor/pubsub/server"
)

const (
	HOST = "localhost"
	PORT = "10905"
	TYPE = "tcp"
)

func main() {

	// Open listener
	l, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("Listening on " + HOST + ":" + PORT)

	// Accept connections

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go server.HandleRequest(conn)
	}
}
