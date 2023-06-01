package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type HealthBar struct {
	x, y            float64
	w, h            float64
	points          int
	maxPoints       int
	healthBarColor  color.RGBA
	healthLostColor color.RGBA
	textFont        font.Face
}

func (h *HealthBar) update(x, y float64, points, maxPoints int) {
	h.x, h.y = x, y
	h.points, h.maxPoints = points, maxPoints
}

func (h *HealthBar) draw(screen *ebiten.Image) {
	x1, y1 := float32(h.x), float32(h.y)
	w1 := float32(h.w * float64(h.points) / float64(h.maxPoints))
	h1 := float32(h.h)
	x2, y2 := float32(h.x)+w1, float32(h.y)
	w2 := float32(h.w * float64(h.maxPoints-h.points) / float64(h.maxPoints))
	h2 := float32(h.h)
	vector.DrawFilledRect(screen, x1, y1, w1, h1, h.healthBarColor, false)
	vector.DrawFilledRect(screen, x2, y2, w2, h2, h.healthLostColor, false)
	healthBarMsg := fmt.Sprintf("%d/%d", h.points, h.maxPoints)
	text.Draw(screen, healthBarMsg, h.textFont, int(h.x), int(h.y), healthBarFontColor)
}
