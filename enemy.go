package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Enemy struct {
	Player
	visibleRange    float64
	dotTargetIndex  int
	hits            []Hit
	greedy          float64
	aggressive      float64
	shootFreq       int
	speedMultiplyer int
}

func (p *Enemy) brain(dots []*Dot, player *Player) {
	if p.greedy > 0.5 && p.aggressive > 0.5 {

	}
	if p.greedy > 0.5 && p.aggressive < 0.5 {

	}
	if p.greedy < 0.5 && p.aggressive > 0.5 {
		if p.detectPlayer(player) {
			p.shootLasers(player)
			p.eatDots(dots)
		} else if p.dotTargetIndex < 0 || p.dotTargetIndex > len(dots)-1 {
			p.searchDots(dots)
		} else if p.dotTargetIndex >= 0 && dots[p.dotTargetIndex] != nil {
			p.moveToTarget(dots)
			p.eatDots(dots)
		}
	}
	if p.greedy < 0.5 && p.aggressive < 0.5 {

	}
}

func (p *Enemy) searchDots(dots []*Dot) {
	if len(dots) > 0 {
		p.dotTargetIndex = rand.Intn(len(dots))
	}
}

func (p *Enemy) update() {
	p.y += p.ySpeed
	p.x += p.xSpeed
	p.healthBar.update(p.x-camX-p.w/2, p.y-(p.h-p.h/3)-camY, p.points, p.maxPoints)
}

func (p *Enemy) draw(screen *ebiten.Image, x float64, y float64, dots []*Dot) {
	// Draw the enemy
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(p.w/2), -float64(p.h/2))
	op.GeoM.Rotate(p.angle)
	op.GeoM.Translate(x, y)
	screen.DrawImage(p.img, op)
	p.healthBar.draw(screen)
	for i := len(p.hits) - 1; i >= 0; i-- {
		if p.hits[i].duration > 0 {
			p.hits[i].update()
			p.hits[i].draw(screen, camX, camY)
		} else {
			p.hits[i] = p.hits[len(p.hits)-1]
			p.hits = p.hits[:len(p.hits)-1]
		}
	}
	ebitenutil.DrawRect(screen, x-p.visibleRange, y-p.visibleRange, p.visibleRange*2, p.visibleRange*2, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x0f})
}

func (p *Enemy) updateLasers() {
	for index, laser := range p.lasers {
		hit := false
		if math.Abs(float64(laser.y+laser.speed*math.Sin(laser.angle))-float64(player.y)) < player.h/2 && math.Abs(float64(laser.x+laser.speed*math.Cos(laser.angle))-float64(player.x)) < player.w/2 {
			player.points -= pointsPerHit
			hit = true
		}
		if hit {
			p.lasers[index] = p.lasers[len(p.lasers)-1]
			p.lasers = p.lasers[:len(p.lasers)-1]
			continue
		}
		laser.update()
	}
}

func (p *Enemy) drawLasers(screen *ebiten.Image, camX float64, camY float64) {
	for _, laser := range p.lasers {
		laser.draw(screen, p.x-camX, p.y-camY)
	}
}

func (p *Enemy) detectPlayer(player *Player) bool {
	if distanceBetweenPoints(player.x, player.y, p.x, p.y) < p.visibleRange {
		return true
	}
	return false
}

func (p *Enemy) moveToTarget(dots []*Dot) {
	p.angle = angleBetweenPoints(p.x, p.y, float64(dots[p.dotTargetIndex].x), float64(dots[p.dotTargetIndex].y))
	p.ySpeed = math.Sin(p.angle) * float64(p.speedMultiplyer)
	p.xSpeed = math.Cos(p.angle) * float64(p.speedMultiplyer)
}

func (p *Enemy) eatDots(dots []*Dot) {
	for dotIndex := range dots {
		if dots[dotIndex] != nil && math.Abs(float64(p.y+p.ySpeed)-float64(dots[dotIndex].y)) < p.h/2 && math.Abs(float64(p.x+p.xSpeed)-float64(dots[dotIndex].x)) < p.w/2 {
			p.points += pointsPerDot
			dots[dotIndex] = nil
			if dotIndex == p.dotTargetIndex {
				p.dotTargetIndex = -1
				p.ySpeed = 0
				p.xSpeed = 0
			}
			if p.points > p.maxPoints {
				p.maxPoints = p.points
			}
		}
	}
}

func (p *Enemy) shootLasers(player *Player) {
	p.angle = angleBetweenPoints(p.x, p.y, player.x, player.y)
	p.ySpeed = math.Sin(p.angle) * float64(p.speedMultiplyer)
	p.xSpeed = math.Cos(p.angle) * float64(p.speedMultiplyer)

	if frameCount%p.shootFreq == 0 {
		p.lasers = append(p.lasers, &Laser{
			x:     p.x,
			y:     p.y,
			angle: p.angle,
			speed: laserSpeed,
		})

		if len(p.lasers) > maxLasers {
			p.lasers[0] = p.lasers[len(p.lasers)-1]
			p.lasers = p.lasers[1:]
		}
	}
}
