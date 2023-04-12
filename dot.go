package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type Dot struct {
	x, y     int
	color    color.RGBA
	msg      string
	textFont font.Face
}

func (d *Dot) draw(screen *ebiten.Image, camX float64, camY float64) {
	text.Draw(screen, d.msg, d.textFont, d.x-int(camX), d.y-int(camY), d.color)
}
