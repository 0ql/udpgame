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
	id      byte
	coord_x uint32
	coord_y uint32
}

func (player *Player) Update() {
	if Win.Pressed(pixelgl.KeyW) {
		player.coord_y += 3
	}
	if Win.Pressed(pixelgl.KeyS) {
		player.coord_y -= 3
	}
	if Win.Pressed(pixelgl.KeyD) {
		player.coord_x += 3
	}
	if Win.Pressed(pixelgl.KeyA) {
		player.coord_x -= 3
	}
}

func (player *Player) Draw() {
	circle := imdraw.New(nil)
	circle.Color = colornames.Black
	circle.Push(pixel.V(float64(player.Coord_x), float64(player.Coord_y)))
	circle.Circle(PLAYER_THICKNESS, 0)
	circle.Draw(Win)
}
