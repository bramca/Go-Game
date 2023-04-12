package main

import (
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten"
)

type Player struct {
	x, y      float64
	w, h      float64
	angle     float64
	lasers    []*Laser
	img       *ebiten.Image
	ySpeed    float64
	xSpeed    float64
	points    int
	maxPoints int
	healthBar HealthBar
}

func (p *Player) update(x, y float64) {
	// Move the player based on the mouse position
	mx, my := ebiten.CursorPosition()
	p.y += p.ySpeed
	p.x += p.xSpeed
	p.healthBar.update(p.x-camX-p.w/2, p.y-(p.h-p.h/3)-camY, p.points, p.maxPoints)
	p.angle = angleBetweenPoints(x, y, float64(mx), float64(my))
}

func (p *Player) draw(screen *ebiten.Image, x float64, y float64) {
	// Draw the player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(p.w/2), -float64(p.h/2))
	op.GeoM.Rotate(p.angle)
	op.GeoM.Translate(x, y)
	screen.DrawImage(p.img, op)
	p.healthBar.draw(screen)
}

func (p *Player) updateLasers() {
	for index, laser := range p.lasers {
		hit := false
		for _, enemy := range enemies {
			if math.Abs(float64(laser.y+laser.speed*math.Sin(laser.angle))-float64(enemy.y)) < enemy.h/2 && math.Abs(float64(laser.x+laser.speed*math.Cos(laser.angle))-float64(enemy.x)) < enemy.w/2 {
				enemy.points -= pointsPerHit
				hit = true
				enemy.hits = append(enemy.hits, Hit{
					Dot: Dot{
						x: int(enemy.x),
						y: int(enemy.y - enemy.h/2),
						color: color.RGBA{
							R: 0xff,
							G: 0xff,
							B: 0xff,
							A: 0xf0,
						},
						msg:      strconv.Itoa(-pointsPerHit),
						textFont: textFont,
					},
					duration: 40,
				})
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
