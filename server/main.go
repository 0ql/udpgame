package main

import (
	"fmt"
	"net"
	"time"
)

type TCPConnection struct {
	addr string
}

type UDPConnection struct {
	addr string
}

func handleTCPConnection(con net.Conn) {
	buf := make([]byte, 100)
	for {
		// read & write timeout
		t := time.Now().Add(50 * time.Millisecond)
		err := con.SetDeadline(t)
		if err != nil {
			fmt.Println("Canceling connection...")
			fmt.Println(err)
			break
		}

		_, err = con.Read(buf)
		if err != nil {
			fmt.Println("Error reading from TCP packet")
			fmt.Println(err)
			continue
		}
	}
}

func createTCPlistener(tcpcon net.Listener) {
	TCPConnections := map[string]TCPConnection{}

	for {
		con, err := tcpcon.Accept()
		if err != nil {
			fmt.Println("TCP Connection Failed")
			continue
		}
		fmt.Println("New TCP Con")

		TCPConnections[con.RemoteAddr().String()] = TCPConnection{
			addr: con.RemoteAddr().String(),
		}

		go handleTCPConnection(con)
	}

}

func handleUDPPackets(con net.PacketConn) {
	buf := make([]byte, 100)
	UDPConnections := map[string]UDPConnection{}

	for {
		_, addr, err := con.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error Reading UDP Packet")
			fmt.Println(err)
			continue
		}

		if _, ok := UDPConnections[addr.String()]; ok {
			// go dostuff()
		} else {
			fmt.Println("New UDP Conn")
			UDPConnections[addr.String()] = UDPConnection{
				addr: addr.String(),
			}
		}
	}
}

func createServer(port string) {
	fmt.Println("Server is starting...")

	tcpLn, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer tcpLn.Close()

	udpcon, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		panic(err)
	}
	defer udpcon.Close()

	go createTCPlistener(tcpLn)
	handleUDPPackets(udpcon)
}

func main() {
	createServer(":8080")
}
