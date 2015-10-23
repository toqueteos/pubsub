package server

import (
	"bufio"
	"io"
	"net"
	"time"

	"github.com/ikanor/pubsub"
	"github.com/satori/go.uuid"
)

type Peer struct {
	pid   uuid.UUID
	start time.Time
	conn  net.Conn
	rw    io.ReadWriter
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		pid:   uuid.NewV4(),
		start: time.Now(),
		conn:  conn,
		rw: bufio.NewReadWriter(
			bufio.NewReader(conn),
			bufio.NewWriter(conn),
		),
	}
}

func (p *Peer) Loop() {
	for {
		message, err := p.readMessage()

		select {
		case <-time.After(pubsub.PeerIdleTimeout):
		}
	}
}

// func (p *Peer) readMessage() (*Message, error) {
// 	return nil, nil
// }

func (p *Peer) Send(payload []byte) error {
	return nil
}
