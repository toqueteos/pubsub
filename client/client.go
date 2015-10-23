package client

import (
	"bytes"
	"io"
	"net"
	"strings"

	"github.com/ikanor/pubsub"
	"github.com/ikanor/pubsub/wire"
)

type SubscribeFn func(client Client, payload []byte)

type Interface interface {
	Subscribe(channel string, fn SubscribeFn) error
	Unsubscribe(channel string) error
	Publish(channel string, payload []byte) error
}

type Client struct {
	addr   string
	conn   net.Conn
	closed bool
}

func New(addr string) (Interface, error) {
	if strings.Index(addr, ":") == -1 {
		addr = addr + ":" + pubsub.DefaultPort
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	cli := &Client{
		addr: addr,
		conn: conn,
	}
	go cli.loop()
	return cli, nil
}

func (c *Client) loop() {
	var buf bytes.Buffer
	for !c.closed {
		io.Copy(&buf, c.conn)
	}
}

// Subscribe tells the server to send us notifications when something is
// published to `channel`.
func (c *Client) Subscribe(channel string, fn SubscribeFn) error {
	var buf bytes.Buffer
	buf.WriteByte(pubsub.CmdSubscribe)
	if err := wire.WriteString(&buf, channel); err != nil {
		return err
	}
	_, err := io.Copy(c.conn, &buf)
	return err
}

// Unsubscribe tells the server to stop sending us notifications when something
// is published to `channel`.
func (c *Client) Unsubscribe(channel string) error {
	var buf bytes.Buffer
	buf.WriteByte(pubsub.CmdUnsubscribe)
	if err := wire.WriteString(&buf, channel); err != nil {
		return err
	}
	_, err := io.Copy(c.conn, &buf)
	return err
}

func (c *Client) Publish(channel string, payload []byte) error {
	panic("client.Publish is not implemented")
	return nil
}

func (c *Client) Stop() {

}
