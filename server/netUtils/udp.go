package netUtils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func NewUDPListener(port string, gameStart time.Time, b *ConBundle) UDPListener {
	GS = gameStart
	udpaddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		panic(err)
	}

	udpLn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		panic(err)
	}

	return UDPListener{
		listener: udpLn,
		bundle:   b,
	}
}

func (udp *UDPListener) SendUDPStatePackets(Hz int) {
	sleeptime := time.Duration(1000 / Hz * int(time.Millisecond))

	b := udp.bundle
	for {
		time.Sleep(sleeptime)
		packet := make([]byte, 1)

		// time since game start
		duration := uint32(time.Since(GS).Milliseconds())
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, duration)
		packet = append(packet, buf...)

		// playercount
		a := uint8(len(b.clients))
		playerAmount := byte(a)
		packet = append(packet, playerAmount)

		for _, v := range b.clients {
			packet = append(packet, v.ID)
			packet = append(packet, v.PS.X...)
			packet = append(packet, v.PS.Y...)
		}

		fmt.Println(packet)
		for _, v := range b.clients {
			_, _, err := udp.listener.WriteMsgUDP(packet, nil, v.addr)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (udp *UDPListener) ReadState(addr net.Addr, packet []byte) {
	b := udp.bundle

	buf := make([]byte, 4)
	for i := 1; i <= 4; i++ {
		buf[i-1] = packet[i]
	}

	var p PlayerState = PlayerState{}

	// check if packet outdated
	if bytes.Compare(buf, b.clients[addr.String()].timestamp) == -1 {
		return
	} else {
		b.clients[addr.String()].timestamp = buf
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

	b.clients[addr.String()].PS = &p
}

func (udp *UDPListener) HandleUDPPackets() {
	buf := make([]byte, 100)

	for {
		_, addr, err := udp.listener.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error Reading UDP Packet")
			fmt.Println(err)
			continue
		}

		udp.ReadState(addr, buf)
	}
}
