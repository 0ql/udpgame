package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type connection struct {
	addr string
}

type state struct {
	timestamp   []byte
	playercount byte
	players     map[byte]Player
	my_id       byte
}

var (
	connections = map[string]connection{}
	gameState   = state{}
)

func New() state {
	gs := state{}
	gs.timestamp = make([]byte, 4)
	gs.playercount = 1
	gs.players = map[byte]Player{}
	gs.my_id = 0
	gs.players[gs.my_id] = Player{
		id:      gs.my_id,
		coord_x: 100,
		coord_y: 100,
	}
	return gs
}

func readToByteArray(buf []byte, leftOffset int, length int) []byte {
	var temp []byte

	for i := leftOffset; i < leftOffset+length; i++ {
		temp = append(temp, buf[i])
	}

	return temp
}

func readToUint64(buf []byte, leftOffset int, length int) uint64 {
	return binary.BigEndian.Uint64(readToByteArray(buf, leftOffset, length))
}

func (self *state) fromBinary(buf []byte) {
	for i := 0; i < 4; i++ {
		self.timestamp = append(self.timestamp, buf[i])
	}

	self.playercount = buf[4]

	for i := 5; i < int(self.playercount); i += 9 {
		var player = Player{}

		player.id = buf[i]
		player.coord_x = readToUint64(buf, i+1, 4)
		player.coord_y = readToUint64(buf, i+2, 4)
		self.players[player.id] = player
	}
}

// func (self *state) toBinary() {
// 	me := self.players[self.my_id]
// 	packet = make([]byte, 13)

// 	for i := 0; i < 4; i++ {
// 		packet
// 	}
// }

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

	buf := make([]byte, 100)

	// helloMsg := []byte("n")

	// con.WriteTo(helloMsg)

	defer con.Close()

	go func() {
		_, addr, err := con.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error receving packet")
		}

		if _, ok := connections[addr.String()]; ok {
			// existing client
		} else {
			// new Client
			connections[addr.String()] = connection{
				addr: addr.String(),
			}
		}
	}()
}
