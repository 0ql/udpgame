package rendering

import (
	"image/color"
	"strings"
	"time"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

func ttfFromBytesMust(b []byte, size float64) font.Face {
	ttf, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(ttf, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}

type typewriter struct {
	regular *text.Text
	text    string

	offset   pixel.Vec
	position pixel.Vec
	move     pixel.Vec
}

func newTypewriter(c color.Color, regular *text.Atlas) *typewriter {
	tw := &typewriter{
		regular: text.New(pixel.ZV, regular),
	}
	tw.regular.Color = c
	return tw
}

func (tw *typewriter) TypeRune(c rune) {
	tw.text += string(c)
	tw.regular.WriteRune(c)
}

func (tw *typewriter) Back() {
	// tw.regular.Dot = tw.regular.Dot.Sub(pixel.V(tw.regular.Atlas().Glyph(' ').Advance, 0))
	if len(tw.text) > 0 {
		// what is this shit??? why isn't there a built-in substring function to remove the last char
		tw.text = strings.TrimSuffix(tw.text, string(tw.text[len(tw.text)-1]))
	}
	tw.regular.Clear()
	tw.regular.WriteString(tw.text)
}

func (tw *typewriter) Offset(off pixel.Vec) {
	tw.offset = tw.offset.Add(off)
}

func (tw *typewriter) Move(vel pixel.Vec) {
	tw.move = vel
}

func (tw *typewriter) Dot() pixel.Vec {
	return tw.regular.Dot
}

func (tw *typewriter) Update(dt float64) {
	tw.position = tw.position.Add(tw.move.Scaled(dt))
}

func (tw *typewriter) Draw(t pixel.Target, m pixel.Matrix) {

	m = pixel.IM.Moved(tw.position.Add(tw.offset)).Chained(m)
	tw.regular.Draw(t, m)
}

type cursorbar struct {
	width  float64
	height float64
	color  color.Color
	pos    pixel.Vec
	imd    *imdraw.IMDraw
}

func CursorBarNew(width, height float64, color color.Color) *cursorbar {
	return &cursorbar{
		width:  width,
		height: height,
		color:  color,
		pos:    pixel.ZV,
		imd:    imdraw.New(nil),
	}
}

func (bar *cursorbar) Update(newpos pixel.Vec) {
	bar.pos = newpos
}

func (bar *cursorbar) Draw(t pixel.Target, m pixel.Matrix) {
	bar.imd.Clear()
	bar.imd.SetMatrix(m)
	bar.imd.Color = bar.color
	bar.imd.Push(bar.pos)
	bar.imd.Push(bar.pos.Add(pixel.V(bar.width, bar.height)))
	bar.imd.Rectangle(0)
	bar.imd.Color = pixel.Alpha(0)
	bar.imd.Draw(t)
}

type textfieldwidget struct {
	regular *text.Atlas
	bgColor color.RGBA
	fgColor color.RGBA
	tw      *typewriter
	bar     *cursorbar
	fps     <-chan time.Time
	last    time.Time
	pos     pixel.Vec
	width   float64
	height  float64
	onEnter func(string)
}

func TextFieldWidgetNew(pos pixel.Vec, width, height float64, onEnter func(string)) textfieldwidget {
	var (
		regular = text.NewAtlas(
			ttfFromBytesMust(goregular.TTF, height),
			text.ASCII, text.RangeTable(unicode.Latin),
		)
		bgColor = colornames.White
		fgColor = colornames.Black

		tw   = newTypewriter(pixel.ToRGBA(fgColor).Scaled(0.9), regular)
		bar  = CursorBarNew(5.0, height, pixel.ToRGBA(fgColor).Scaled(0.9))
		fps  = time.Tick(time.Second / 120)
		last = time.Now()
	)

	widget := textfieldwidget{
		regular: regular,
		bgColor: bgColor,
		fgColor: fgColor,
		tw:      tw,
		bar:     bar,
		fps:     fps,
		last:    last,
		pos:     pos,
		width:   width,
		height:  height,
		onEnter: onEnter,
	}

	return widget
}

func (widget *textfieldwidget) Update(win *pixelgl.Window) {
	for _, c := range win.Typed() {
		widget.tw.TypeRune(c)
	}

	if win.JustPressed(pixelgl.KeyTab) || win.Repeated(pixelgl.KeyTab) {
		widget.tw.regular.WriteRune('\t')
	}
	if win.JustPressed(pixelgl.KeyBackspace) || win.Repeated(pixelgl.KeyBackspace) {
		widget.tw.Back()
	}

	if win.JustPressed(pixelgl.KeyEnter) {
		widget.onEnter(widget.tw.text)
	}

	dt := time.Since(widget.last).Seconds()
	widget.last = time.Now()

	widget.tw.Update(dt)
	widget.bar.Update(widget.tw.Dot())
}

func (widget *textfieldwidget) Draw(win *pixelgl.Window) {
	m := pixel.IM.Moved(widget.pos)
	widget.tw.Draw(win, m)
	widget.bar.Draw(win, m)
	<-widget.fps
}
