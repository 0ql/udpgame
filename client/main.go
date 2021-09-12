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
	frames     = 0
	second     = time.Tick(time.Second)
	PlayerName = ""
	TcpCon     n.TCPCon
	widget     = rendering.TextFieldWidgetNew(pixel.V(32, 32), 200.0, 64.0, func(text string) {
		if PlayerName == "" {
			PlayerName = text
			TcpCon.SendConnectRequestPacket(PlayerName)
		}
	})
	label = rendering.LabelNew("Please enter your nickname below", pixel.V(32, 200), 200.0, 20.0)
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "udpgame",
		Bounds: pixel.R(0, 0, 600, 500),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)

		if PlayerName != "" {
			rendering.StateMutex.Lock()
			for key := range rendering.GS.Players {
				player := rendering.GS.Players[key]
				if player.Id == rendering.GS.My_id {
					player.Update(win)
				}
				player.Draw(win)
			}
			rendering.StateMutex.Unlock()

			frames++
			select {
			case <-second:
				win.SetTitle(fmt.Sprintf("%s | FPS: %d | TCP %d↓ %d↑ | UDP %d↓ %d↑", cfg.Title, frames, n.TCPPPSDOWN, n.TCPPPSUP, n.UDPPPSDOWN, n.UDPPPSUP))
				frames = 0
				n.TCPPPSDOWN = 0
				n.TCPPPSUP = 0
				n.UDPPPSDOWN = 0
				n.UDPPPSUP = 0
			default:
			}
		} else {
			widget.Update(win)
			widget.Draw(win)
			label.Draw(win)
		}

		win.Update()
	}
}

func main() {
	tcpConnection, err := n.NewTCPConn("localhost:8080")
	if err != nil {
		panic(err)
	}

	go tcpConnection.SendStayAlivePackets()
	go tcpConnection.ListenPackets()

	TcpCon = tcpConnection

	pixelgl.Run(run)
}
