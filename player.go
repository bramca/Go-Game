package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Player struct {
	x, y         float64
	w, h         float64
	angle        float64
	lasers       []*Laser
	img          *ebiten.Image
	ySpeed       float64
	xSpeed       float64
	points       int
	maxPoints    int
	healthBar    HealthBar
	score        int
	fireRate     int
	laserSpeed   float64
	speed        float64
	acceleration float64
	damage       int
}

func (p *Player) update(x, y float64, dots []*Dot) {
	// Move the player based on the mouse position
	mx, my := ebiten.CursorPosition()
	p.y += p.ySpeed
	p.x += p.xSpeed
	p.healthBar.update(p.x-camX-p.w/2, p.y-(p.h-p.h/3)-camY, p.points, p.maxPoints)
	p.angle = angleBetweenPoints(x, y, float64(mx), float64(my))
	for dotIndex := range dots {
		if !dots[dotIndex].eaten && dots[dotIndex].duration > 0 && (distanceBetweenPoints(p.x+p.xSpeed, p.y+p.ySpeed, float64(dots[dotIndex].x), float64(dots[dotIndex].y)) < p.w*0.8 || distanceBetweenPoints(p.x+p.xSpeed, p.y+p.ySpeed, float64(dots[dotIndex].x+len(dots[dotIndex].msg)), float64(dots[dotIndex].y)) < p.w) {
			p.points += pointsPerDot
			dots[dotIndex].hits = append(dots[dotIndex].hits, Hit{
				Dot: Dot{
					x:        dots[dotIndex].x,
					y:        dots[dotIndex].y,
					color:    dotHitColor,
					msg:      "+" + strconv.Itoa(pointsPerDot),
					textFont: hitTextFont,
				},
				duration: 2 * framesPerSecond / 3,
			})
			dots[dotIndex].eaten = true
			if p.points > p.maxPoints {
				p.maxPoints = p.points
			}
		}
	}
}

func (p *Player) drawScore(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("Score: %d", p.score), scoreTextFont, scoreFontSize, scoreFontSize+10, scoreColor)
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
	for index := len(p.lasers) - 1; index >= 0; index-- {
		hit := false
		for _, enemy := range enemies {
			if !enemy.dead && math.Abs(float64(p.lasers[index].y+p.lasers[index].speed*math.Sin(p.lasers[index].angle))-float64(enemy.y)) < enemy.h/2 && math.Abs(float64(p.lasers[index].x+p.lasers[index].speed*math.Cos(p.lasers[index].angle))-float64(enemy.x)) < enemy.w/2 {
				enemy.points -= p.lasers[index].damage
				hit = true
				enemy.hits = append(enemy.hits, Hit{
					Dot: Dot{
						x:        int(enemy.x),
						y:        int(enemy.y - enemy.h/2),
						color:    damageColor,
						msg:      strconv.Itoa(-p.lasers[index].damage),
						textFont: hitTextFont,
					},
					duration: 2 * framesPerSecond / 3,
				})
			}
		}
		for _, lootBox := range lootBoxes {
			if !lootBox.broken && lootBox.duration > 0 && math.Abs(float64(p.lasers[index].y+p.lasers[index].speed*math.Sin(p.lasers[index].angle))-float64(lootBox.y)) < lootBox.h/2 && math.Abs(float64(p.lasers[index].x+p.lasers[index].speed*math.Cos(p.lasers[index].angle))-float64(lootBox.x)) < lootBox.w/2 {
				lootBox.hitpoints -= p.lasers[index].damage
				hit = true
				lootBox.hits = append(lootBox.hits, Hit{
					Dot: Dot{
						x: int(lootBox.x),
						y: int(lootBox.y - lootBox.h/2),
						color: color.RGBA{
							R: 0xff,
							G: 0xff,
							B: 0xff,
							A: 0xf0,
						},
						msg:      strconv.Itoa(-p.lasers[index].damage),
						textFont: hitTextFont,
					},
					duration: 2 * framesPerSecond / 3,
				})
			}
		}
		if hit {
			p.lasers[index] = p.lasers[len(p.lasers)-1]
			p.lasers = p.lasers[:len(p.lasers)-1]
			continue
		}
		p.lasers[index].update()
	}
}

func (p *Player) drawLasers(screen *ebiten.Image, camX float64, camY float64) {
	for index := len(p.lasers) - 1; index >= 0; index-- {
		if p.lasers[index].duration < 0 {
			p.lasers[index] = p.lasers[len(p.lasers)-1]
			p.lasers = p.lasers[:len(p.lasers)-1]
			continue
		}
		p.lasers[index].draw(screen, camX, camY)
	}
}
