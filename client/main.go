package main

import (
	"fmt"

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

	gameState = StateNew()

	for !Win.Closed() {
		Win.Clear(colornames.Whitesmoke)

		fmt.Println(gameState)
		for key := range gameState.players {
			player := gameState.players[key]
			player.Draw()
		}

		Win.Update()
	}
}

func main() {
	// go StartConnection("localhost:8080")
	pixelgl.Run(run)
}
