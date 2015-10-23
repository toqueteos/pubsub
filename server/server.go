package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/ikanor/pubsub"
	"github.com/ikanor/pubsub/client"

	"github.com/satori/go.uuid"
	"gopkg.in/inconshreveable/log15.v2"
)

type Server struct {
	sid      uuid.UUID
	addr     string
	listener net.Listener
	errors   chan error
	clients  map[string]*client.Client
	// subscribers struct {
	// 	byId map[string]*Client
	// 	byCh map[string]*[]Client
	// }

	peerMessages chan PeerMessage
}

func New(addr string) (*Server, error) {
	if strings.Index(addr, ":") == -1 {
		addr = fmt.Sprintf("%s:%s", addr, pubsub.DefaultPort)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		sid:      uuid.NewV4(),
		addr:     addr,
		listener: listener,
		errors:   make(chan error),
		clients:  make(map[string]*client.Client),
	}, nil
}

func (s *Server) Start() {
	go s.pingPeers()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.errors <- err
			continue
		}

		peer := NewPeer(conn)
		go peer.Loop()
		go s.addPeer(peer)
	}
}

func (s *Server) loop() {
	var ping = time.Tick(30 * time.Second)
	for {
		select {
		case <-ping:
			log15.Info("Pinging all clients...")
			go s.pingPeers()
		}
	}
}

func (s Server) pingPeers() {

}

type PeerMessage struct {
	Command byte
	Peer    *Peer
}

func (s *Server) addPeer(addr net.Addr, peer *Peer) {
	s.peerMessages <- PeerMessage{CmdPeerAdd, peer}
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

type CmdInfo struct {
	Conn    net.Conn
	Command string
	Channel string
	Payload []byte
}

func (s *Server) handleRequest(conn net.Conn) {
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		command, err := r.ReadByte()
		if err != nil {
			s.errors <- err
			continue
		}
		channel, err := readSizeAndBlob(r).String()
		if err != nil {
			s.errors <- err
			continue
		}
		payload, err := readSizeAndBlob(r).Bytes()
		if err != nil {
			s.errors <- err
			continue
		}

		s.processCommand(&CmdInfo{
			Conn:    conn,
			Command: command,
			Channel: channel,
			Payload: payload,
		})
	}
}

func readSizeAndBlob(r io.Reader) (*bytes.Buffer, error) {
	var (
		buf  bytes.Buffer
		size int64
	)
	err := binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	buf.Reset()
	io.CopyN(&buf, r, size)
	return buf, nil
}

func (s *Server) ping(conn net.Conn) error {
	_, err := fmt.Fprintf(conn, "%x-pubsub-ping-%d", s.sid, time.Now().Unix())
}
