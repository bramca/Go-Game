package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type HealthBar struct {
	x, y      float64
	w, h      float64
	points    int
	maxPoints int
}

func (h *HealthBar) update(x, y float64, points, maxPoints int) {
	h.x, h.y = x, y
	h.points, h.maxPoints = points, maxPoints
}

func (h *HealthBar) draw(screen *ebiten.Image) {
	x1, y1 := h.x, h.y
	w1 := h.w * float64(h.points) / float64(h.maxPoints)
	h1 := h.h
	x2, y2 := h.x+w1, h.y
	w2 := h.w * float64(h.maxPoints-h.points) / float64(h.maxPoints)
	h2 := h.h
	ebitenutil.DrawRect(screen, x1, y1, w1, h1, color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xf0})
	ebitenutil.DrawRect(screen, x2, y2, w2, h2, color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xf0})
}
