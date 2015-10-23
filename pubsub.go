package pubsub

import "time"

const DefaultPort = "10905"

var (
	PeerIdleTimeout = 5 * time.Second
)

const (
	CmdError byte = iota
	CmdRegister
	CmdPublish
	CmdSubscribe
	CmdUnsubscribe
	CmdStart
	CmdEnd
	CmdPing
	CmdPong
	CmdPeerAdd
	CmdPeerDel
)

const (
	WireByte byte = iota
	WireString
	WireSlice
)
