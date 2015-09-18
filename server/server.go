package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
)

type Server struct {
	addr     string
	listener net.Listener
	errors   chan error
}

func NewServer(addr string) (*Server, error) {
	if strings.Index(addr, ":") == -1 {
		addr = fmt.Sprintf("%s:%s", addr, DefaultPort)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		addr:     addr,
		listener: listener,
		errors:   make(chan error, 32),
	}, nil
}

func (s *Server) Start() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.errors <- err
			continue
		}
		go s.handleRequest(conn)
	}
}

func (s *Server) ForEachError(fn func(error)) {
	for err := range s.errors {
		fn(err)
	}
}

func (s *Server) Stop() {
	s.listener.Close()
	close(s.errors)
}

func (s *Server) handleRequest(conn net.Conn) {
	r := bufio.NewReadWriter(bufio.NewReader(conn), nil)

	cmd, err := r.ReadString(' ')
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

// protoRead reads binary blobs with the following format:
//
//    COMMAND ' ' CHANNEL ' ' SIZE PAYLOAD
func (s *Server) protoRead(conn net.Conn) {
	r := bufio.NewReader(conn)
	rw := bufio.NewReadWriter(r, nil)

	var payloadSize int64
	err := binary.Read(rw, binary.BigEndian, &payloadSize)

	var buf bytes.Buffer
	io.CopyN(&buf, rw, payloadSize)

	body := buf.Bytes()
	commandIndex := bytes.Index(body, []byte(" "))
	channelIndex := bytes.Index(body[commandIndex+1:], []byte(" "))

	blobs := bytes.SplitN(body, []byte(" "), 2)
	command := string(blobs[0])
	channel := string(blobs[1])
	payload := blobs[2]

	reply, err := processCommand(command, channel, payload)
	if err != nil {
		fmt.Println("Error processing command:", err)
		return
	}

	conn.Write([]byte(reply))
	// Close the connection when you're done with it.
	conn.Close()
}

func processCommand(cmd string, msg string) (string, error) {
	switch cmd {
	case "PUB":
		// do publish
	case "SUB":
		// do subscribe
	case "UNS":
		// do unsubscribe
	default:
		return "NACK\n", fmt.Errorf("Unknown command %q", cmd)
	}
	return "ACK\n", nil
}
