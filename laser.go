package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Laser struct {
	x, y     float64
	angle    float64
	speed    float64
	size     float64
	color    color.RGBA
	duration int
	damage   int
}

func (l *Laser) update() {
	l.x += l.speed * math.Cos(l.angle)
	l.y += l.speed * math.Sin(l.angle)
	l.duration -= 1
}

func (l *Laser) draw(screen *ebiten.Image, x float64, y float64) {
	// Draw the laser line
	laserX := l.x - x
	laserY := l.y - y
	x1, y1 := laserX, laserY
	x2, y2 := laserX+l.size*math.Cos(l.angle), laserY+l.size*math.Sin(l.angle)
	vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 3, l.color, false)
}
