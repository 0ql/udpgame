package main

import (
	n "client/networking"
	"client/rendering"
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	frames = 0
	second = time.Tick(time.Second)
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "udpgame",
		Bounds: pixel.R(0, 0, 500, 500),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)

		for key := range rendering.GS.Players {
			player := rendering.GS.Players[key]
			if player.Id == rendering.GS.My_id {
				player.Update(win)
			}
			player.Draw(win)
		}

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	tcpConnection, err := n.NewTCPConn("localhost:8080")
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
