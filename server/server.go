package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"

	. "github.com/ikanor/pubsub/conf"
)

type Server struct {
	addr        string
	listener    net.Listener
	errors      chan error
	subscribers struct {
		byId map[string]*Client
		byCh map[string]*[]Client
	}
}

type Client struct {
	addr string
	uuid string
}

func New(addr string) (*Server, error) {
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
		clients:  make(map[string]*Client),
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

// protoRead reads binary blobs with the following format:
//
//    COMMAND ' ' CHANNEL ' ' SIZE PAYLOAD
func (s *Server) handleRequest(conn net.Conn) {
	r := bufio.NewReader(conn)
	rw := bufio.NewReadWriter(r, nil)

	var payloadSize int64
	err := binary.Read(rw, binary.BigEndian, &payloadSize)

	var buf bytes.Buffer
	io.CopyN(&buf, rw, payloadSize)

	body := buf.Bytes()

	blobs := bytes.SplitN(body, []byte(" "), 2)
	command := string(blobs[0])
	channel := string(blobs[1])
	payload := blobs[2]

	reply, err := s.processCommand(command, channel, payload)
	if err != nil {
		fmt.Println("Error processing command:", err)
		return
	}

	conn.Write([]byte(reply))
	// Close the connection when you're done with it.
	conn.Close()
}

func (s *Server) processCommand(command string, client *Client, channel string, payload []byte) (string, error) {

	// We need to register the clients on the first command
	// s.subscribers.byId[uuid] = &Client{uuid: uuid.New()}

	switch command {
	case PUB:
		s.publish(channel, payload)
	case SUB:
		s.subscribe(client, channel)
	case UNS:
		// do unsubscribe
	default:
		return NACK, fmt.Errorf("Unknown command %q", command)
	}
	return ACK, nil
}

func (s *Server) publish(channel string, payload string) {
	// publish
}

func (s *Server) subscribe(client Client, channel string) {
	// subscribe
}
