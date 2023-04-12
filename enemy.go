package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Enemy struct {
	Player
	visibleRange   float64
	dotTargetIndex int
	hits           []Hit
}

func (p *Enemy) searchDots(screen *ebiten.Image, dots []*Dot) {
	if len(dots) > 0 {
		p.dotTargetIndex = rand.Intn(len(dots))
	}
}

func (p *Enemy) update(dots []*Dot) {
	if p.dotTargetIndex >= 0 && dots[p.dotTargetIndex] != nil {
		p.angle = angleBetweenPoints(p.x, p.y, float64(dots[p.dotTargetIndex].x), float64(dots[p.dotTargetIndex].y))
		p.ySpeed = math.Sin(p.angle) * 2
		p.xSpeed = math.Cos(p.angle) * 2
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
	if p.dotTargetIndex >= 0 && p.dotTargetIndex < len(dots) && dots[p.dotTargetIndex] != nil {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\ndot.x, dot.y: %d, %d\nangleBetween: %02f", dots[p.dotTargetIndex].x, dots[p.dotTargetIndex].y, p.angle))
		ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\nenemy.xSpeed, enemy.ySpeed, enemy.points: %02f, %02f, %d", p.xSpeed, p.ySpeed, p.points))
		ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\ndistanceToY, distanceToX: %02f, %02f", math.Abs(float64(p.y+p.ySpeed)-float64(dots[p.dotTargetIndex].y)), math.Abs(float64(p.x+p.xSpeed)-float64(dots[p.dotTargetIndex].x))))
	}
}

func (p *Enemy) updateLasers() {
	for _, laser := range p.lasers {
		laser.update()
	}
}

func (p *Enemy) drawLasers(screen *ebiten.Image, camX float64, camY float64) {
	for _, laser := range p.lasers {
		laser.draw(screen, p.x-camX, p.y-camY)
	}
}

func (p *Enemy) detectPlayer(screen *ebiten.Image, player *Player) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("enemy.x, enemy.y: %02f, %02f\nplayer.x, player.y: %02f, %02f", p.x, p.y, player.x, player.y))
}
