package rendering

import (
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

type labelwidget struct {
	regular *text.Text
	pos     pixel.Vec
	width   float64
	height  float64
}

func LabelNew(labelText string, pos pixel.Vec, width, height float64) labelwidget {
	regular := text.NewAtlas(
		ttfFromBytesMust(goregular.TTF, height),
		text.ASCII, text.RangeTable(unicode.Latin),
	)

	widget := labelwidget{
		regular: text.New(pos, regular),
		pos:     pos,
		width:   width,
		height:  height,
	}

	widget.regular.Color = colornames.Black
	widget.regular.WriteString(labelText)

	return widget
}

func (widget *labelwidget) Draw(t pixel.Target) {
	widget.regular.Draw(t, pixel.IM.Moved(widget.pos))
}
