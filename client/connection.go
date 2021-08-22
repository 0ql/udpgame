package main

import (
	"fmt"
	"net"
	"time"
)

type connection struct {
	addr       string
	con        *net.UDPConn
	start_time time.Time
}

var (
	gameConnection connection
)

func StartConnection(serverIP string) {
	fmt.Println("Connecting to Server")

	raddr, err := net.ResolveUDPAddr("udp", serverIP)
	if err != nil {
		panic(err)
	}

	con, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		panic(err)
	}

	gameConnection = connection{
		addr:       raddr.String(),
		con:        con,
		start_time: time.Now(),
	}

	// helloMsg := []byte("n")
	// con.WriteTo(helloMsg)

	defer con.Close()

	go func() {
		buf := make([]byte, 100)
		for {
			_, _, err := con.ReadFrom(buf)
			if err != nil {
				fmt.Println("Error receving packet")
				continue
			}

			fromBinary(buf)
		}
	}()
}

func fromBinary(buf []byte) {
	if buf[0] == STATE_PACKET_ID {
		// init state
		// state.fromBinary(buf)
	}
}
