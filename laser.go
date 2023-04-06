package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Laser struct {
	x, y  float64
	angle float64
	speed float64
}

func (l *Laser) update() {
	l.x += l.speed * math.Cos(l.angle)
	l.y += l.speed * math.Sin(l.angle)
}

func (l *Laser) draw(screen *ebiten.Image, camX float64, camY float64) {
	// Draw the laser line
	laserX := l.x - camX
	laserY := l.y - camY
	x1, y1 := laserX, laserY
	x2, y2 := laserX+10*math.Cos(l.angle), laserY+10*math.Sin(l.angle)
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{R: 255, G: 0, B: 0, A: 255})
}
