package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

type Player struct {
	x, y   float64
	w, h   float64
	angle  float64
	lasers []*Laser
	img    *ebiten.Image
	ySpeed float64
	xSpeed float64
	points int
}

func (p *Player) update(x, y float64) {
	// Move the player based on the mouse position
	mx, my := ebiten.CursorPosition()
	p.y += p.ySpeed
	p.x += p.xSpeed
	p.angle = angleBetweenPoints(x, y, float64(mx), float64(my))
}

func (p *Player) draw(screen *ebiten.Image, x float64, y float64) {
	// Draw the player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(p.w/2), -float64(p.h/2))
	op.GeoM.Rotate(p.angle)
	op.GeoM.Translate(x, y)
	screen.DrawImage(p.img, op)
}

func (p *Player) updateLasers() {
	for index, laser := range p.lasers {
		hit := false
		for _, enemy := range enemies {
			if math.Abs(float64(laser.y+laser.speed*math.Sin(laser.angle))-float64(enemy.y)) < enemy.h/2 && math.Abs(float64(laser.x+laser.speed*math.Cos(laser.angle))-float64(enemy.x)) < enemy.w/2 {
				enemy.points -= 1
				hit = true
			}
		}
		if hit {
			p.lasers[index] = p.lasers[len(p.lasers)-1]
			p.lasers = p.lasers[:len(p.lasers)-1]
			continue
		}
		laser.update()
	}
}

func (p *Player) drawLasers(screen *ebiten.Image, camX float64, camY float64) {
	for _, laser := range p.lasers {
		laser.draw(screen, camX, camY)
	}
}
