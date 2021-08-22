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

func (self *state) fromBinary(buf []byte) {
	self.timestamp = BufChunkToByteArray(buf, 1, 4)

	self.playercount = buf[5]

	for i := 6; i < int(self.playercount); i += 9 {
		var player = Player{}

		player.id = buf[i]
		player.coord_x = BufChunkToUint64(buf, i+1, 4)
		player.coord_y = BufChunkToUint64(buf, i+2, 4)
		self.players[player.id] = player
	}
}

func (self *state) toBinary() []byte {
	my_player := gameState.players[gameState.my_id]

	packet := make([]byte, 14)

	packet[0] = STATE_PACKET_ID

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(time.Now().Sub(gameConnection.start_time).Milliseconds()))
	InsertBufChunkInBuf(packet, b)

	packet[6] = self.my_id
	InsertBufChunkInBuf(packet, my_player.coord_x, 7)
	InsertBufChunkInBuf(packet, my_player.coord_y, 11)

	return packet
}
