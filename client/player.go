package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Player struct {
	id      byte
	coord_x uint64
	coord_y uint64
}

func (player *Player) Draw() {
	fmt.Println("draw")
	circle := imdraw.New(nil)
	circle.Color = colornames.Black
	circle.Push(pixel.V(float64(player.coord_x), float64(player.coord_y)))
	circle.Circle(PLAYER_THICKNESS, 0)
	circle.Draw(Win)

}
