package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	Player
	visibleRange       float64
	shootRange         float64
	dotTargetIndex     int
	hits               []Hit
	greedy             float64
	aggressive         float64
	shootFreq          int
	speedMultiplyer    int
	movementPrediction float64
	dead               bool
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
		} else if p.dotTargetIndex >= 0 && !dots[p.dotTargetIndex].eaten {
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
	p.damage = p.maxPoints / 10
}

func (p *Enemy) draw(screen *ebiten.Image, x float64, y float64, dots []*Dot) {
	// Draw the enemy
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(p.w/2), -float64(p.h/2))
	op.GeoM.Rotate(p.angle)
	op.GeoM.Translate(x, y)
	screen.DrawImage(p.img, op)
	p.healthBar.draw(screen)
}

func (p *Enemy) drawHits(screen *ebiten.Image) {
	for i := len(p.hits) - 1; i >= 0; i-- {
		if p.hits[i].duration > 0 {
			p.hits[i].update()
			p.hits[i].draw(screen, camX, camY)
		} else {
			p.hits[i] = p.hits[len(p.hits)-1]
			p.hits = p.hits[:len(p.hits)-1]
		}
	}
}

func (p *Enemy) updateLasers() {
	for index := len(p.lasers) - 1; index >= 0; index-- {
		hit := false
		if math.Abs(float64(p.lasers[index].y+p.lasers[index].speed*math.Sin(p.lasers[index].angle))-float64(player.y)) < player.h/2 && math.Abs(float64(p.lasers[index].x+p.lasers[index].speed*math.Cos(p.lasers[index].angle))-float64(player.x)) < player.w/2 {
			player.points -= p.lasers[index].damage
			hit = true
		}
		if hit {
			p.lasers[index] = p.lasers[len(p.lasers)-1]
			p.lasers = p.lasers[:len(p.lasers)-1]
			continue
		}
		p.lasers[index].update()
	}
}

func (p *Enemy) drawLasers(screen *ebiten.Image, camX float64, camY float64) {
	for index := len(p.lasers) - 1; index >= 0; index-- {
		if p.lasers[index].duration < 0 {
			p.lasers[index] = p.lasers[len(p.lasers)-1]
			p.lasers = p.lasers[1:]
			continue
		}
		p.lasers[index].draw(screen, p.x-camX, p.y-camY)
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
		if !dots[dotIndex].eaten && math.Abs(float64(p.y+p.ySpeed)-float64(dots[dotIndex].y)) < p.h/2 && math.Abs(float64(p.x+p.xSpeed)-float64(dots[dotIndex].x)) < p.w/2 {
			p.points += pointsPerDot
			dots[dotIndex].eaten = true
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
	p.angle = angleBetweenPoints(p.x, p.y, player.x+player.xSpeed*p.movementPrediction, player.y+player.ySpeed*p.movementPrediction)
	p.ySpeed = math.Sin(p.angle) * float64(p.speedMultiplyer)
	p.xSpeed = math.Cos(p.angle) * float64(p.speedMultiplyer)

	if distanceBetweenPoints(p.x, p.y, player.x, player.y) <= p.shootRange {
		p.ySpeed = 0
		p.xSpeed = 0
	}

	if frameCount%p.shootFreq == 0 {
		p.lasers = append(p.lasers, &Laser{
			x:        p.x,
			y:        p.y,
			angle:    p.angle,
			speed:    laserSpeed,
			color:    enemyLaserColor,
			duration: laserDuration,
			size:     laserSize,
			damage:   p.damage,
		})
	}
}
