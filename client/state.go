package main

import (
	"encoding/binary"
	"time"
)

type state struct {
	timestamp   []byte
	playercount byte
	players     map[byte]Player
	my_id       byte
}

var (
	gameState = state{}
)

func StateNew() state {
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

func (own_state *state) UpdateFromPacket(packetDecoder PacketDecoder) {
	own_state.timestamp = packetDecoder.ExtractData(4)
	own_state.playercount = packetDecoder.ExtractByte()

	for i := 6; i < int(own_state.playercount); i += 9 {
		player := Player{}

		player.id = packetDecoder.ExtractByte()
		player.coord_x = binary.BigEndian.Uint64(packetDecoder.ExtractData(4))
		player.coord_x = binary.BigEndian.Uint64(packetDecoder.ExtractData(4))

		own_state.players[player.id] = player
	}
}

func (own_state *state) ToPacket(start_time time.Time) []byte {
	my_player := gameState.players[gameState.my_id]

	packet := PacketBuilderNew(UDP_STATE_PACKET)

	packet.add_data(Uint64ToByteArray(uint64(time.Since(start_time).Milliseconds())))

	packet.add_byte(my_player.id)

	packet.add_data(Uint64ToByteArray(my_player.coord_x))
	packet.add_data(Uint64ToByteArray(my_player.coord_y))

	return packet.build()
}
