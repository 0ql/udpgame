package rendering

import (
	"encoding/binary"
	"fmt"
	"time"
)

type state struct {
	Timestamp   []byte
	playercount byte
	Players     map[byte]Player
	My_id       byte
}

var (
	GS = state{
		Players: map[byte]Player{},
	}
)

func (s *state) UpdateFromInitialStatePacket(packet []byte) {
	var p Player
	p.Id = packet[1]
	s.My_id = packet[1]

	buf := make([]byte, 4)
	for i := 2; i <= 5; i++ {
		buf[i-2] = packet[i]
	}

	s.Timestamp = buf

	s.playercount = packet[6]
	for i := 7; i <= 10; i++ {
		buf[i-7] = packet[i]
		fmt.Println(i)
	}
	p.Coord_x = binary.BigEndian.Uint32(buf)

	for i := 11; i <= 14; i++ {
		buf[i-11] = packet[i]
	}
	p.Coord_y = binary.BigEndian.Uint32(buf)

	s.Players[packet[1]] = p
}

func (s *state) UpdateFromPacket(packet []byte) {
	s.playercount = packet[5]

	for i := 0; i < int(s.playercount); i++ {
		coord_x := make([]byte, 4)
		for j := i; j < i+4; j++ {
			coord_x[j-i*8] = packet[j]
		}

		coord_y := make([]byte, 4)
		for j := i + 4; j < i+8; j++ {
			coord_y[j-i*8] = packet[j]
		}
		s.Players[packet[i*8]] = Player{
			Id:      packet[i*8],
			Coord_x: binary.BigEndian.Uint32(coord_x),
			Coord_y: binary.BigEndian.Uint32(coord_y),
		}
	}
}

func (s *state) ToPacket(st time.Time) []byte {
	packet := make([]byte, 1)

	duration := uint32(time.Since(st).Milliseconds())
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, duration)
	packet = append(packet, buf...)

	binary.BigEndian.PutUint32(buf, s.Players[s.My_id].Coord_x)
	packet = append(packet, buf...)
	binary.BigEndian.PutUint32(buf, s.Players[s.My_id].Coord_y)
	packet = append(packet, buf...)

	return packet
}
