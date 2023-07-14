package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type RubberDuck struct {
	Player
	visibleRange       float64
	fleeRange          float64
	dotTargetIndex     int
	hits               []Hit
	speedMultiplyer    int
	movementPrediction float64
	reward             string
	rewardGiven        bool
	dead               bool
	fleeing            bool
}

func (p *RubberDuck) giveReward() {
	switch p.reward {
	case "Exploding Lasers":
		player.gun = "Exploding Lasers"
	case "Double Lasers":
		player.gun = "Double Lasers"
	case "Piercing Lasers":
		player.gun = "Piercing Lasers"
	case "Homing Lasers":
		player.gun = "Homing Lasers"
	case "Shotgun":
		player.gun = "Shotgun"
	}
	p.rewardGiven = true
}

func (p *RubberDuck) brain(dots []*Dot, player *Player) {
	if p.fleeing {
		p.moveAwayFromTarget(player)
		if p.checkFleeDistance(player) {
			p.fleeing = false
		}
		p.eatDots(dots)
	} else if p.detectPlayer(player) {
		p.moveAwayFromTarget(player)
		p.fleeing = true
		p.eatDots(dots)
	} else if p.dotTargetIndex < 0 || p.dotTargetIndex > len(dots)-1 {
		p.searchDots(dots)
	} else if p.dotTargetIndex >= 0 && !dots[p.dotTargetIndex].eaten {
		p.moveToTarget(dots)
		p.eatDots(dots)
	}
}

func (p *RubberDuck) searchDots(dots []*Dot) {
	if len(dots) > 0 {
		p.dotTargetIndex = rand.Intn(len(dots))
	}
}

func (p *RubberDuck) update() {
	p.y += p.ySpeed
	p.x += p.xSpeed
	p.healthBar.update(p.x-camX-p.w/2, p.y-(p.h-p.h/3)-camY, p.points, p.maxPoints)
}

func (p *RubberDuck) draw(screen *ebiten.Image, x float64, y float64) {
	// Draw the rubberDuck
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(p.w/2), -float64(p.h/2))
	op.GeoM.Rotate(p.angle)
	op.GeoM.Translate(x, y)
	screen.DrawImage(p.img, op)
	p.healthBar.draw(screen)
}

func (p *RubberDuck) drawHits(screen *ebiten.Image) {
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

func (p *RubberDuck) detectPlayer(player *Player) bool {
	return distanceBetweenPoints(player.x, player.y, p.x, p.y) < p.visibleRange
}

func (p *RubberDuck) checkFleeDistance(player *Player) bool {
	return distanceBetweenPoints(player.x, player.y, p.x, p.y) > p.fleeRange
}

func (p *RubberDuck) moveToTarget(dots []*Dot) {
	p.angle = angleBetweenPoints(p.x, p.y, float64(dots[p.dotTargetIndex].x), float64(dots[p.dotTargetIndex].y))
	p.ySpeed = math.Sin(p.angle) * float64(p.speedMultiplyer)
	p.xSpeed = math.Cos(p.angle) * float64(p.speedMultiplyer)
}

func (p *RubberDuck) moveAwayFromTarget(player *Player) {
	if !p.fleeing || rand.Float64() < 0.05 {
		min := -math.Pi / 4
		max := math.Pi / 4
		r := math.Pi + (min + rand.Float64()*(max-min))
		p.angle = angleBetweenPoints(p.x, p.y, float64(player.x), float64(player.y)) + r
	}
	p.ySpeed = math.Sin(p.angle) * float64(p.speedMultiplyer)
	p.xSpeed = math.Cos(p.angle) * float64(p.speedMultiplyer)
}

func (p *RubberDuck) eatDots(dots []*Dot) {
	for dotIndex := range dots {
		if !dots[dotIndex].eaten && dots[dotIndex].duration > 0 && (distanceBetweenPoints(p.x+p.xSpeed, p.y+p.ySpeed, float64(dots[dotIndex].x), float64(dots[dotIndex].y)) < p.w*0.8 || distanceBetweenPoints(p.x+p.xSpeed, p.y+p.ySpeed, float64(dots[dotIndex].x+len(dots[dotIndex].msg)), float64(dots[dotIndex].y)) < p.w) {
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
