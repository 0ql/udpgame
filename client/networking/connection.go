package networking

import (
	"bytes"
	r "client/rendering"
	"errors"
	"fmt"
	"net"
	"time"
)

var (
	serverAddr net.Addr
	TCPPPSDOWN = 0
	TCPPPSUP   = 0
	UDPPPSDOWN = 0
	UDPPPSUP   = 0
)

type UDPCon struct {
	raddr      net.UDPAddr
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

func NewUDPConn(localAddr net.Addr, raddr net.UDPAddr) (*UDPCon, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", localAddr.String())
	if err != nil {
		panic(err)
	}

	udpLn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("LOCAL UDP ADDR: %s \n", udpLn.LocalAddr().String())

	return &UDPCon{
		raddr:      raddr,
		con:        udpLn,
		start_time: time.Now(),
		timestamp:  r.GS.Timestamp,
	}, nil
}

func (udp *UDPCon) sendStatePackets(Hz int) error {
	dur := time.Duration(1000 / Hz * int(time.Millisecond))
	for {
		time.Sleep(dur)

		p := r.GS.ToPacket(udp.start_time)
		_, _, err := udp.con.WriteMsgUDP(p, nil, &udp.raddr)
		UDPPPSUP++
		if err != nil {
			panic(err)
		}
	}
}

// blocking
func (udp *UDPCon) ListenPackets() error {
	buf := make([]byte, 100)
	t := make([]byte, 4)
	for {
		_, addr, err := (*udp.con).ReadFrom(buf)
		UDPPPSDOWN++
		if err != nil {
			panic(err)
		}

		if addr.String() != serverAddr.String() {
			continue
		}

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
	serverAddr = con.RemoteAddr()

	if err != nil {
		return TCPCon{}, err
	}

	return TCPCon{
		con:        con,
		addr:       address,
		start_time: time.Now(),
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
	TCPPPSUP++
	if err != nil {
		return err
	}

	return nil
}

func (tcp *TCPCon) sendPlayerListRequestPacket() error {
	packet := make([]byte, 0)

	packet = append(packet, byte(2))
	_, err := tcp.con.Write(packet)
	TCPPPSUP++
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
		fmt.Printf("TCP SAL to: %s \n", tcp.addr)
		_, err := tcp.con.Write(packet)
		TCPPPSUP++
		if err != nil {
			tcp.errorChannel <- err
			break
		}
	}
}

func (tcp *TCPCon) startUDP() {
	udpRAdddr, err := net.ResolveUDPAddr("udp", tcp.con.RemoteAddr().String())
	if err != nil {
		panic(err)
	}
	udpCon, err := NewUDPConn(tcp.con.LocalAddr(), *udpRAdddr)
	if err != nil {
		panic(err)
	}
	tcp.udpCon = udpCon

	fmt.Println("Waiting for UDP Packet...")
	go tcp.udpCon.ListenPackets()
	go tcp.udpCon.sendStatePackets(30)
}

func (tcp *TCPCon) ListenPackets() error {
	buf := make([]byte, 100)

	for {
		err := tcp.con.SetDeadline(time.Now().Add(time.Second))
		if err != nil {
			return err
		}

		_, err = tcp.con.Read(buf)
		TCPPPSDOWN++
		if err != nil {
			fmt.Println("Error receiving TCP packet")
			return err
		}

		switch buf[0] {
		case 0:
			fmt.Printf("TCPConnection Confirmed: Player ID: %d \n", buf[1])
			r.GS.UpdateFromInitialStatePacket(buf)
			go tcp.sendStayAlivePackets()
			// tcp.SendPlayerListRequestPacket()
			tcp.startUDP()
		case 1:
			continue
		case 2:
			// save playerlist
			continue
		case 4:
			r.GS.RemovePlayer(buf[1])
		}
	}

}
