package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	STATE_PACKET_TIMEOUT = "33ms" // not exactly 30 Hz, maybe rather 25Hz (40ms)
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

func (tcp *TCPCon) waitForTCPConnectPacket() {
	buf := make([]byte, 100)
	tcp.con.Read(buf)

	pd := PacketDecoderNew(buf)

	if pd.GetPacketType() != TCP_CONNECT_PACKET {
		fmt.Println("wrong packet type")
	}

	gameState.my_id = pd.ExtractByte()
}

func (tcp *TCPCon) sendChunkRequestPacket() {
	pb := PacketBuilderNew(TCP_CHUNK_REQUEST_PACKET)
	pb.add_data(Uint64ToByteArray(0))
	pb.add_data(Uint64ToByteArray(0))

	tcp.con.Write(pb.build())
}

func (tcp *TCPCon) sendPlayerlistRequestPacket() {
	pb := PacketBuilderNew(TCP_PLAYERLIST_REQUEST_PACKET)

	tcp.con.Write(pb.build())
}

func (tcp *TCPCon) waitForTCPPlayerlistPacket() {
	buf := make([]byte, 100)
	tcp.con.Read(buf)

	pd := PacketDecoderNew(buf)

	if pd.GetPacketType() != TCP_PLAYERLIST_PACKET {
		fmt.Println("wrong packet type")
	}

	gameState.playercount = pd.ExtractByte()

	for i := 0; i < int(gameState.playercount); i++ {
		player := Player{
			id:      pd.ExtractByte(),
			coord_x: uint64(binary.BigEndian.Uint64(pd.ExtractData(4))),
			coord_y: uint64(binary.BigEndian.Uint64(pd.ExtractData(4))),
		}

		gameState.players[player.id] = player
	}
}

func (udp *UDPCon) handleUDPPackets() {
	buf := make([]byte, 100)
	for {
		_, _, err := udp.con.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error receving packet")
			continue
		}

		packetDecoder := PacketDecoderNew(buf)

		if packetDecoder.GetPacketType() == UDP_STATE_PACKET {
			gameState.UpdateFromPacket(packetDecoder)
		}
	}
}

func (udp *UDPCon) sendStatePackets() {
	timeout, err := time.ParseDuration(STATE_PACKET_TIMEOUT)
	if err != nil {
		panic(err)
	}

	for {
		udp.con.Write(gameState.ToPacket(udp.start_time))
		time.Sleep(timeout)
	}
}
