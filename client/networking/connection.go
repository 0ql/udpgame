package networking

import (
	"bytes"
	r "client/rendering"
	"errors"
	"fmt"
	"net"
	"time"
)

type UDPCon struct {
	addr       string
	con        *net.UDPConn
	start_time time.Time
	timestamp  []byte
}

type TCPCon struct {
	addr         string
	con          *net.TCPConn
	udpCon       *UDPCon
	start_time   time.Time
	errorChannel chan error
}

func (udp *UDPCon) sendStatePackets(Hz int) error {
	dur := time.Duration(1000 / Hz)
	for {
		time.Sleep(dur)

		_, err := udp.con.Write(r.GS.ToPacket(udp.start_time))
		if err != nil {
			return err
		}
	}
}

func (udp *UDPCon) ListenPackets() error {
	buf := make([]byte, 100)
	t := make([]byte, 4)
	for {
		_, err := udp.con.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println(buf)

		for i := 1; i <= 4; i++ {
			t = append(t, buf[i])
		}

		// check if packet is old
		if bytes.Compare(udp.timestamp, t) == 1 {
			continue
		} else {
			udp.timestamp = t
		}

		r.GS.UpdateFromPacket(buf)
	}
}

func NewTCPConn(address string) (TCPCon, error) {
	tcpaddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
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

func NewUDPConn(serverAddress string) (UDPCon, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		return UDPCon{}, nil
	}

	con, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		return UDPCon{}, nil
	}

	return UDPCon{
		con:        con,
		addr:       serverAddress,
		start_time: time.Now(),
		timestamp:  r.GS.Timestamp,
	}, nil
}

func (tcp *TCPCon) SendConnectRequestPacket(playerName string) error {
	fmt.Println("Connecting to Server...")
	packet := make([]byte, 1)

	// don't have to set packet type because []byte by default Zeros

	if len(playerName) > 8 {
		return errors.New("playername too long")
	}
	name := []byte("player")
	packet = append(packet, name...)

	_, err := tcp.con.Write(packet)
	if err != nil {
		return err
	}

	return nil
}

func (tcp *TCPCon) SendPlayerListRequestPacket() error {
	packet := make([]byte, 0)

	packet = append(packet, byte(2))
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
			tcp.errorChannel <- err
			break
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

		switch buf[0] {
		case 0:
			fmt.Println(buf)
			r.GS.UpdateFromInitialStatePacket(buf)
			go tcp.sendStayAlivePackets()
			tcp.SendPlayerListRequestPacket()
			udpCon, err := NewUDPConn(tcp.addr)
			fmt.Println("Starting UDP Connection...")
			if err != nil {
				panic(err)
			}
			tcp.udpCon = &udpCon
			go tcp.udpCon.sendStatePackets(30)
		case 1:
			continue
		case 2:
			// save playerlist
			continue
		}
	}

}
