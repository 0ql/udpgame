package rendering

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	PLAYER_THICKNESS float64 = 10
)

type Player struct {
	Id      byte
	Coord_x uint32
	Coord_y uint32
}

func (player *Player) Update(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyW) {
		player.Coord_y += 3
	}
	if win.Pressed(pixelgl.KeyS) {
		player.Coord_y -= 3
	}
	if win.Pressed(pixelgl.KeyD) {
		player.Coord_x += 3
	}
	if win.Pressed(pixelgl.KeyA) {
		player.Coord_x -= 3
	}
}

func (player *Player) Draw(win *pixelgl.Window) {
	circle := imdraw.New(nil)
	circle.Color = colornames.Black
	circle.Push(pixel.V(float64(player.Coord_x), float64(player.Coord_y)))
	circle.Circle(PLAYER_THICKNESS, 0)
	circle.Draw(win)
}
