package server

import (
	"bufio"
	"fmt"
	"net"
)

func HandleRequest(conn net.Conn) {
	r := bufio.NewReadWriter(bufio.NewReader(conn), nil)

	cmd, err := r.ReadString(byte(' '))
	if err != nil {
		fmt.Println("Error reading command:", err.Error())
		return
	}
	cmd = cmd[:len(cmd)-1]
	raw, _, err := r.ReadLine()
	if err != nil {
		fmt.Println("Error reading meesage:", err.Error())
		return
	}
	msg := string(raw)

	fmt.Printf("Read: %s - %s\n", cmd, msg)

	reply, err := processCommand(cmd, msg)
	if err != nil {
		fmt.Println("Error processing command:", err.Error())
		return
	}

	conn.Write([]byte(reply))
	// Close the connection when you're done with it.
	conn.Close()
}

func processCommand(cmd string, msg string) (string, error) {
	switch cmd {
	case "PUBLISH":
		// do publish
	case "SUBSCRIBE":
		// do subscribe
	case "UNSUBSCRIBE":
		// do unsubscribe
	default:
		return "NACK\n", fmt.Errorf("Unknown command %q", cmd)
	}
	return "ACK\n", nil
}
