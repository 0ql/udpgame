package main

import (
	"client/networking"
	"client/rendering"

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

	rendering.GameState = rendering.StateNew()

	for !Win.Closed() {
		Win.Clear(colornames.Whitesmoke)

		for key := range rendering.GameState.Players {
			player := rendering.GameState.Players[key]
			player.Draw(Win)
		}

		Win.Update()
	}
}

func main() {
	tcpconnection, err := networking.NewTCPConn("localhost:8080")
	if err != nil {
		panic(err)
	}

	err = tcpconnection.SendConnectRequestPacket("hans")
	if err != nil {
		panic(err)
	}

	err = tcpconnection.ListenPackets()
	if err != nil {
		panic(err)
	}
	// go StartConnection("localhost:8080")
	// pixelgl.Run(run)
}
