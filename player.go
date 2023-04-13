package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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

func (p *Player) update(x, y float64, dots []*Dot) {
	// Move the player based on the mouse position
	mx, my := ebiten.CursorPosition()
	p.y += p.ySpeed
	p.x += p.xSpeed
	p.healthBar.update(p.x-camX-p.w/2, p.y-(p.h-p.h/3)-camY, p.points, p.maxPoints)
	p.angle = angleBetweenPoints(x, y, float64(mx), float64(my))
	for dotIndex := range dots {
		if dots[dotIndex] != nil && math.Abs(float64(p.y+p.ySpeed)-float64(dots[dotIndex].y)) < p.h/2 && math.Abs(float64(p.x+p.xSpeed)-float64(dots[dotIndex].x)) < p.w/2 {
			p.points += pointsPerDot
			dots[dotIndex] = nil
			if p.points > p.maxPoints {
				p.maxPoints = p.points
			}
		}
	}
}

func (p *Player) draw(screen *ebiten.Image, x float64, y float64) {
	// Draw the player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(p.w/2), -float64(p.h/2))
	op.GeoM.Rotate(p.angle)
	op.GeoM.Translate(x, y)
	screen.DrawImage(p.img, op)
	p.healthBar.draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\nplayer.x, player.y: %02f, %02f", p.x, p.y))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\nplayer.xSpeed, player.ySpeed, player.points: %02f, %02f, %d", p.xSpeed, p.ySpeed, p.points))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\ndistance player - enemy: %02f", distanceBetweenPoints(p.x, p.y, enemies[0].x, enemies[0].y)))
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
