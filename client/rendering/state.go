package rendering

import (
	"client/util"
	"encoding/binary"
	"time"
)

type state struct {
	Timestamp   []byte
	Playercount byte
	Players     map[byte]Player
	My_id       byte
}

var (
	GameState = state{}
)

func StateNew() state {
	gs := state{}
	gs.Timestamp = make([]byte, 4)
	gs.Playercount = 1
	gs.Players = map[byte]Player{}
	gs.My_id = 0
	gs.Players[gs.My_id] = Player{
		Id:      gs.My_id,
		Coord_x: 100,
		Coord_y: 100,
	}
	return gs
}

func (own_state *state) UpdateFromPacket(packetDecoder util.PacketDecoder) {
	own_state.Timestamp = packetDecoder.ExtractData(4)
	own_state.Playercount = packetDecoder.ExtractByte()

	for i := 6; i < int(own_state.Playercount); i += 9 {
		Player := Player{}

		Player.Id = packetDecoder.ExtractByte()
		Player.Coord_x = binary.BigEndian.Uint64(packetDecoder.ExtractData(4))
		Player.Coord_x = binary.BigEndian.Uint64(packetDecoder.ExtractData(4))

		own_state.Players[Player.Id] = Player
	}
}

func (own_state *state) ToPacket(start_time time.Time) []byte {
	my_Player := GameState.Players[GameState.My_id]

	packet := util.PacketBuilderNew(util.UDP_STATE_PACKET)

	packet.AddData(util.Uint64ToByteArray(uint64(time.Since(start_time).Milliseconds())))

	packet.AddByte(my_Player.Id)

	packet.AddData(util.Uint64ToByteArray(my_Player.Coord_x))
	packet.AddData(util.Uint64ToByteArray(my_Player.Coord_y))

	return packet.Build()
}
