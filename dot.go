package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Dot struct {
	x, y        int
	color       color.RGBA
	msg         string
	textFont    *text.GoXFace
	drawOptions *text.DrawOptions
	hits        []Hit
	eaten       bool
	duration    int
}

func (d *Dot) update() {
	d.duration -= 1
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
	d.drawOptions.DrawImageOptions.GeoM.Translate(float64(d.x-int(camX)), float64(d.y-int(camY)))
	text.Draw(screen, d.msg, d.textFont, d.drawOptions)
	d.drawOptions.DrawImageOptions.GeoM.Reset()
}
