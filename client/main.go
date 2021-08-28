package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	Win *pixelgl.Window
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "udpgame",
		Bounds: pixel.R(0, 0, 1000, 500),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	Win = win

	if err != nil {
		panic(err)
	}

	for !Win.Closed() {
		Win.Clear(colornames.Whitesmoke)

		for key := range GS.players {
			player := GS.players[key]
			if player.id == GS.my_id {
				player.Update()
			}
			player.Draw()
		}

		Win.Update()
	}
}

func main() {
	tcpConnection, err := NewTCPConn("localhost:8080")
	if err != nil {
		panic(err)
	}

	go tcpConnection.ListenPackets()

	err = tcpConnection.SendConnectRequestPacket("hans")
	if err != nil {
		panic(err)
	}

	pixelgl.Run(run)
}
