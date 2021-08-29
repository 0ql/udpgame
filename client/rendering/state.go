package rendering

import (
	"encoding/binary"
	"fmt"
	"time"
)

type state struct {
	Timestamp   []byte
	playercount byte
	Players     map[byte]*Player
	My_id       byte
}

var (
	GS = state{
		Players: map[byte]*Player{},
	}
)

func (s *state) UpdateFromInitialStatePacket(packet []byte) {
	fmt.Println("Initial RSU")
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

	s.Players[packet[1]] = &p
}

func (s *state) UpdateFromPacket(packet []byte) {
	fmt.Println("RSU")
	s.playercount = packet[5]

	for i := 0; i < int(s.playercount); i++ {
		for j := 6 + i*7; j <= 6+i*7+7; j++ {
			coord_x := make([]byte, 4)

			for k := 0; k <= 3; k++ {
				coord_x[k] = packet[j+1+k]
			}

			coord_y := make([]byte, 4)
			for k := 0; k <= 3; k++ {
				coord_y[k] = packet[j+4+k]
			}

			s.Players[packet[i*7+6]] = &Player{
				Id:      packet[i*7+6],
				Coord_x: binary.BigEndian.Uint32(coord_x),
				Coord_y: binary.BigEndian.Uint32(coord_y),
			}
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
