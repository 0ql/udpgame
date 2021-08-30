package netUtils

import (
	"net"
	"time"
)

// gamestart
var GS time.Time

type ConBundle struct {
	tcp                  map[string]*TCPConnection
	clients              map[string]*Client
	removeConnectionChan chan string
	MAXPPS               uint
}

type TCPConnection struct {
	addr         string
	Con          net.Conn
	ID           byte // same as player ID
	pps          *uint
	parentBundle *ConBundle
}

type PlayerState struct {
	X []byte
	Y []byte
}

type Client struct {
	addr      *net.UDPAddr
	ID        byte
	PS        *PlayerState
	timestamp []byte
}

type UDPListener struct {
	listener *net.UDPConn
	bundle   *ConBundle
}
