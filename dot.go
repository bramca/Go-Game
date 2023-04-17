package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Dot struct {
	x, y     int
	color    color.RGBA
	msg      string
	textFont font.Face
	hits     []Hit
	eaten    bool
}

func (d *Dot) drawHits(screen *ebiten.Image, camX float64, camY float64) {
	for i := len(d.hits) - 1; i >= 0; i-- {
		if d.hits[i].duration > 0 {
			d.hits[i].update()
			d.hits[i].draw(screen, camX, camY)
		} else {
			d.hits[i] = d.hits[len(d.hits)-1]
			d.hits = d.hits[:len(d.hits)-1]
		}
	}
}

func (d *Dot) draw(screen *ebiten.Image, camX float64, camY float64) {
	text.Draw(screen, d.msg, d.textFont, d.x-int(camX), d.y-int(camY), d.color)
}
