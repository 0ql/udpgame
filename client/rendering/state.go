package rendering

import (
	"client/util"
	"encoding/binary"
	"sync"
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

	StateMutex = sync.Mutex{}
)

func (s *state) UpdateFromInitialStatePacket(packet []byte) {
	StateMutex.Lock()

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
	}
	p.Coord_x = binary.BigEndian.Uint32(buf)

	for i := 11; i <= 14; i++ {
		buf[i-11] = packet[i]
	}
	p.Coord_y = binary.BigEndian.Uint32(buf)

	s.Players[packet[1]] = &p

	StateMutex.Unlock()
}

func (s *state) UpdateFromPacket(packet []byte) {
	StateMutex.Lock()

	d := util.NewPacketDecoder(packet)
	d.SetIndex(5)
	s.playercount = d.ExtractByte()

	for i := 0; i < int(s.playercount); i++ {
		id := d.ExtractByte()
		coord_x := d.ExtractData(4)
		coord_y := d.ExtractData(4)
		s.Players[id] = &Player{
			Id:      id,
			Coord_x: binary.BigEndian.Uint32(coord_x),
			Coord_y: binary.BigEndian.Uint32(coord_y),
		}
	}

	StateMutex.Unlock()
}

func (s *state) ToPacket(st time.Time) []byte {
	StateMutex.Lock()

	packet := make([]byte, 1)

	duration := uint32(time.Since(st).Milliseconds())
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, duration)
	packet = append(packet, buf...)

	binary.BigEndian.PutUint32(buf, s.Players[s.My_id].Coord_x)
	packet = append(packet, buf...)
	binary.BigEndian.PutUint32(buf, s.Players[s.My_id].Coord_y)
	packet = append(packet, buf...)

	StateMutex.Unlock()

	return packet
}

func (s *state) RemovePlayer(id byte) {
	StateMutex.Lock()

	delete(s.Players, id)

	StateMutex.Unlock()
}
