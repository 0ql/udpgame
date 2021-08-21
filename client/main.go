package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var Win *pixelgl.Window

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

		Win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
