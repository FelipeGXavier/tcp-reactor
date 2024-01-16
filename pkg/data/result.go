package data

import "net"

type Result struct {
	Socket net.Conn
	Data string
}