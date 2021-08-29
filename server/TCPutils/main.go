package TCPutils

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

var gs time.Time

type TCPConBundle struct {
	Connections          map[string]*TCPConnection
	removeConnectionChan chan string
	MAXPPS               uint
}

type TCPConnection struct {
	addr         string
	Con          net.Conn
	ID           byte // same as player ID
	pps          *uint
	parentBundle *TCPConBundle
}

func (tcp *TCPConnection) createConnectRequestResponsePacket() []byte {
	packet := make([]byte, 1)

	packet = append(packet, tcp.ID)

	duration := uint32(time.Since(gs).Milliseconds())
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, duration)

	packet = append(packet, byte(len(tcp.parentBundle.Connections)))

	buf = make([]byte, 4)
	packet = append(packet, buf...)
	packet = append(packet, buf...)
	return packet
}

func (tcp *TCPConnection) createStayAliveConfirmationPacket() []byte {
	packet := make([]byte, 0)
	return append(packet, byte(3))
}

func (tcp *TCPConnection) handleTCPConnection(kill chan string) {
	buf := make([]byte, 100)
	var pack []byte

	for {
		// read & write timeout
		t := time.Now().Add(1000 * time.Millisecond)
		err := tcp.Con.SetDeadline(t)
		if err != nil {
			fmt.Println(err)
			kill <- tcp.addr
			break
		}

		_, err = tcp.Con.Read(buf)
		if err != nil {
			kill <- tcp.addr
			fmt.Println(err)
			break
		}

		*tcp.pps++

		if buf[0] == 0 {
			pack = tcp.createConnectRequestResponsePacket()
		} else if buf[0] == 4 {
			fmt.Printf("TCP SAL to: %s at %d PPS \n", tcp.addr, *tcp.pps)
			pack = tcp.createStayAliveConfirmationPacket()
		}

		_, err = tcp.Con.Write(pack)
		if err != nil {
			kill <- tcp.addr
			break
		}
	}
}

func NewTCPConBundle(maxpps uint, gameStart time.Time) TCPConBundle {
	gs = gameStart
	return TCPConBundle{
		Connections:          map[string]*TCPConnection{},
		removeConnectionChan: make(chan string),
		MAXPPS:               maxpps,
	}
}

// blocking
func (bundle *TCPConBundle) ConnectionRemover() {
	for {
		address := <-bundle.removeConnectionChan
		delete(bundle.Connections, address)
		fmt.Printf("Current Connection Count: %d\n", len(bundle.Connections))
	}
}

// blocking
func (bundle *TCPConBundle) checkPpsAndReset() {
	for {
		time.Sleep(time.Second)
		for k, v := range bundle.Connections {
			if *v.pps > bundle.MAXPPS {
				bundle.Connections[k].Con.Close()
				bundle.removeConnectionChan <- v.addr
			} else {
				*v.pps = 0
			}
		}
	}
}

func (bundle *TCPConBundle) findAvailableID() (byte, bool) {
	for i := byte(0); i < 255; i++ {
		found := true
		for _, v := range bundle.Connections {
			if i == v.ID {
				found = false
			}
		}
		if found {
			return i, true
		}
	}
	return byte(0), false
}

// blocking
func (bundle *TCPConBundle) CreateTCPlistener(port string) error {
	tcpLn, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	defer tcpLn.Close()
	go bundle.checkPpsAndReset()

	for {
		con, err := tcpLn.Accept()
		if err != nil {
			fmt.Println("TCP Connection Failed")
			continue
		}
		fmt.Println("New TCP Con")

		pps := uint(0)
		if id, found := bundle.findAvailableID(); found {

			bundle.Connections[con.RemoteAddr().String()] = &TCPConnection{
				addr:         con.RemoteAddr().String(),
				Con:          con,
				ID:           id,
				pps:          &pps,
				parentBundle: bundle,
			}

			c := bundle.Connections[con.RemoteAddr().String()]
			go c.handleTCPConnection(bundle.removeConnectionChan)
		}

	}
}
