package client

import (
	"bufio"
	"errors"
	"net"

	. "github.com/ikanor/pubsub/conf"
)

var ErrNack = errors.New("The server responded with a NACK")

type Client interface {
	Subscribe(string, func(string)) (*Subscription, error)
	Unsubscribe(*Subscription) error
	Publish(string) error
}

type Subscription struct {
	name   string
	action func(string)
}

type client struct {
	addr string
	subs map[string][]*Subscription
	conn net.Conn
	buff *bufio.ReadWriter
}

func NewClient(addr string) Client {
	return &client{addr: addr}
}

func (c *client) connect() (err error) {
	c.conn, err = net.Dial("tcp", c.addr)
	if err != nil {
		return
	}
	c.buff = bufio.NewReadWriter(bufio.NewReader(c.conn), nil)
	return
}

func (c *client) Subscribe(ch string, action func(string)) (*Subscription, error) {
	if c.conn == nil {
		_ = c.connect()
	}
	if len(c.subs[ch]) == 0 {
		c.buff.Write([]byte(SUB + " " + ch))
		msg, err := c.buff.ReadString('\x00')
		if err != nil {
			return nil, err
		}
		if msg == NACK {
			return nil, ErrNack
		}
	}
	s := Subscription{
		name:   ch,
		action: action,
	}
	c.subs[ch] = append(c.subs[ch], &s)
	return &s, nil
}

func (c *client) Unsubscribe(s *Subscription) error {
	n := c.subs[s.name]
	for i := range n {
		if n[i] != s {
			continue
		}
		n[i] = n[len(n)-1]
		c.subs[s.name] = n[:len(n)-1]
		break
	}
	if len(c.subs[s.name]) == 0 {
		c.buff.Write([]byte(UNS + " " + s.name))
		msg, err := c.buff.ReadString('\x00')
		if err != nil {
			return err
		}
		if msg == NACK {
			return ErrNack
		}
	}
	return nil

}

func (c *client) Publish(msg string) error {
	return nil
}
