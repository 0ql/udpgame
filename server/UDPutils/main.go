package UDPutils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"server/TCPutils"
	"time"
)

var gs time.Time

type PlayerState struct {
	X         []byte
	Y         []byte
	Timestamp []byte
}

type Client struct {
	addr net.Addr
	ID   byte
	PS   *PlayerState
}

type UDPListener struct {
	listener  net.PacketConn
	tcpBundle *TCPutils.TCPConBundle
	players   map[string]*Client
}

func NewUDPListener(port string, gameStart time.Time, tcpBundle *TCPutils.TCPConBundle) UDPListener {
	gs = gameStart

	udpLn, err := net.ListenPacket("udp", port)
	if err != nil {
		panic(err)
	}

	return UDPListener{
		listener:  udpLn,
		tcpBundle: tcpBundle,
		players:   map[string]*Client{},
	}
}

func (udp *UDPListener) SendUDPStatePackets(Hz int) {
	sleeptime := time.Duration(1000 / Hz * int(time.Millisecond))

	for {
		time.Sleep(sleeptime)
		packet := make([]byte, 1)

		// time since game start
		duration := uint32(time.Since(gs).Milliseconds())
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, duration)
		packet = append(packet, buf...)

		// playercount
		a := uint8(len(udp.players))
		playerAmount := byte(a)
		packet = append(packet, playerAmount)

		for _, v := range (*udp).players {
			packet = append(packet, v.PS.Timestamp...)
			packet = append(packet, v.PS.X...)
			packet = append(packet, v.PS.Y...)
		}

		for _, v := range udp.players {
			udp.listener.WriteTo(packet, v.addr)
		}
	}
}

func (udp *UDPListener) ReadState(addr net.Addr, packet []byte) {
	buf := make([]byte, 4)
	for i := 1; i <= 4; i++ {
		buf[i-1] = packet[i]
	}

	var p PlayerState = PlayerState{}

	// check if packet outdated
	if bytes.Compare(buf, udp.players[addr.String()].PS.Timestamp) == -1 {
		return
	} else {
		p.Timestamp = buf
	}

	buf = make([]byte, 4)
	for i := 5; i <= 8; i++ {
		buf[i-5] = packet[i]
	}
	p.X = buf

	buf = make([]byte, 4)
	for i := 9; i <= 12; i++ {
		buf[i-9] = packet[i]
	}
	p.Y = buf

	udp.players[addr.String()].PS = &p
}

func (udp *UDPListener) HandleUDPPackets() {
	buf := make([]byte, 100)

	var id uint8 = 0
	for {
		_, addr, err := udp.listener.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error Reading UDP Packet")
			fmt.Println(err)
			continue
		}

		if _, ok := udp.players[addr.String()]; !ok {
			fmt.Println("NEW UDP CON")
			// new UDP Con
			udp.players[addr.String()] = &Client{
				addr: addr,
				ID:   byte(id),
				PS: &PlayerState{
					Timestamp: make([]byte, 0),
					X:         make([]byte, 4),
					Y:         make([]byte, 4),
				},
			}
			id++
			fmt.Println(byte(id))
		}
		udp.ReadState(addr, buf)
	}
}
