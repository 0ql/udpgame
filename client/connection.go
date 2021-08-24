package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type UDPCon struct {
	addr       string
	con        *net.UDPConn
	start_time time.Time
}

type TCPCon struct {
	addr       string
	con        *net.TCPConn
	start_time time.Time
}

func NewTCPConn(address string) (TCPCon, error) {
	tcpaddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println("sdf")
		return TCPCon{}, err
	}

	con, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		return TCPCon{}, err
	}

	return TCPCon{
		con:        con,
		addr:       address,
		start_time: time.Now(),
	}, nil
}

func NewUDPConn(serverURL string, port string) (UDPCon, error) {
	udpaddr, err := net.ResolveUDPAddr(serverURL, port)
	if err != nil {
		return UDPCon{}, nil
	}

	con, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		return UDPCon{}, nil
	}

	return UDPCon{
		con:        con,
		addr:       serverURL,
		start_time: time.Now(),
	}, nil
}

func (tcp *TCPCon) SendConnectRequestPacket(playerName string) error {
	packet := make([]byte, 1)

	// don't have to set packet type because []byte by default Zeros

	if len(playerName) > 8 {
		return errors.New("Playername too long")
	}
	name := []byte("player")
	packet = append(packet, name...)

	_, err := tcp.con.Write(packet)
	if err != nil {
		return err
	}

	return nil
}

// blocking
func (tcp *TCPCon) sendStayAlivePackets() {
	packet := make([]byte, 0)

	packet = append(packet, byte(4))
	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Sending Stay Alive Packet to: %s \n", tcp.addr)
		_, err := tcp.con.Write(packet)
		if err != nil {
			panic(err)
		}
	}
}

func (tcp *TCPCon) ListenPackets() error {
	buf := make([]byte, 100)

	for {
		err := tcp.con.SetDeadline(time.Now().Add(time.Second))
		if err != nil {
			return err
		}

		_, err = tcp.con.Read(buf)
		if err != nil {
			fmt.Println("Error receiving TCP packet")
			return err
		}

		if buf[0] == 0 {
			go tcp.sendStayAlivePackets()
		}
	}

}
