package UDPutils

import (
	"encoding/binary"
	"fmt"
	"net"
	"server/TCPutils"
	"time"
)

var gs time.Time

type UDPListener struct {
	listener  net.PacketConn
	tcpBundle *TCPutils.TCPConBundle
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
	}
}

func (udp *UDPListener) SendUDPStatePackets(Hz int) {
	sleeptime := time.Duration(1000 / Hz * int(time.Millisecond))

	for {
		packet := make([]byte, 1)
		time.Sleep(sleeptime)

		duration := uint32(time.Since(gs).Milliseconds())
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, duration)
		packet = append(packet, buf...)

		a := uint8(len(udp.tcpBundle.Connections))
		playerAmount := byte(a)

		packet = append(packet, playerAmount)

		fmt.Println(packet)
		for _, v := range udp.tcpBundle.Connections {
			udp.listener.WriteTo(packet, v.Con.RemoteAddr())
		}
	}
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

		if _, ok := udp.tcpBundle.Connections[addr.String()]; ok {
			// go doStuff()
		}
	}
}
