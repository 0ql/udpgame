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
	Coord_x uint64
	Coord_y uint64
}

func (player *Player) Draw(win *pixelgl.Window) {
	circle := imdraw.New(nil)
	circle.Color = colornames.Black
	circle.Push(pixel.V(float64(player.Coord_x), float64(player.Coord_y)))
	circle.Circle(PLAYER_THICKNESS, 0)
	circle.Draw(win)

}
