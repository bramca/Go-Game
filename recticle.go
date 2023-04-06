package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Recticle struct {
	size int
}

func (r *Recticle) draw(screen *ebiten.Image) {
	mx, my := ebiten.CursorPosition()
	x1, y1 := float64(mx), float64(my - r.size - 1)
	x2, y2 := float64(mx), float64(my - 1)
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.White)
	x3, y3 := float64(mx - r.size - 1), float64(my)
	x4, y4 := float64(mx - 1), float64(my)
	ebitenutil.DrawLine(screen, x3, y3, x4, y4, color.White)
	x5, y5 := float64(mx), float64(my + r.size + 1)
	x6, y6 := float64(mx), float64(my + 1)
	ebitenutil.DrawLine(screen, x5, y5, x6, y6, color.White)
	x7, y7 := float64(mx + r.size + 1), float64(my)
	x8, y8 := float64(mx + 1), float64(my)
	ebitenutil.DrawLine(screen, x7, y7, x8, y8, color.White)
}
