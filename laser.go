package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Laser struct {
	x, y              float64
	angle             float64
	speed             float64
	size              float64
	color             color.RGBA
	duration          int
	damage            int
	homing            bool
	homingTargetIndex int
	homingTargetType  string
	homingRange       float64
	piercing          bool
	exploding         bool
}

func (l *Laser) update() {
	if l.homing {
		if l.homingTargetIndex < 0 && len(enemies) > 0 && len(rubberDucks) > 0 {
			l.searchEnemy()
		}
		switch l.homingTargetType {
		case "Enemy":
			if l.homingTargetIndex < 0 || l.homingTargetIndex > len(enemies)-1 {
				l.searchEnemy()
			}
			if l.homingTargetIndex >= 0 && l.homingTargetIndex < len(enemies) && !enemies[l.homingTargetIndex].dead && distanceBetweenPoints(enemies[l.homingTargetIndex].x, enemies[l.homingTargetIndex].y, l.x, l.y) < l.homingRange {
				l.angle = angleBetweenPoints(l.x, l.y, enemies[l.homingTargetIndex].x, enemies[l.homingTargetIndex].y)
			} else {
				l.searchEnemy()
			}
		case "RubberDuck":
			if l.homingTargetIndex < 0 || l.homingTargetIndex > len(rubberDucks)-1 {
				l.searchEnemy()
			}
			if l.homingTargetIndex >= 0 && l.homingTargetIndex < len(rubberDucks) && !rubberDucks[l.homingTargetIndex].dead && distanceBetweenPoints(rubberDucks[l.homingTargetIndex].x, rubberDucks[l.homingTargetIndex].y, l.x, l.y) < l.homingRange {
				l.angle = angleBetweenPoints(l.x, l.y, rubberDucks[l.homingTargetIndex].x, rubberDucks[l.homingTargetIndex].y)
			} else {
				l.searchEnemy()
			}
		}
	}
	l.x += l.speed * math.Cos(l.angle)
	l.y += l.speed * math.Sin(l.angle)
	l.duration -= 1
}

func (l *Laser) searchEnemy() {
	l.homingTargetIndex = 0
	minDist := distanceBetweenPoints(l.x, l.y, enemies[0].x, enemies[0].y)
	if len(enemies) > 0 {
		for index, enemy := range enemies {
			d := distanceBetweenPoints(l.x, l.y, enemy.x, enemy.y)
			if d < minDist {
				minDist = d
				l.homingTargetIndex = index
				l.homingTargetType = "Enemy"
			}
		}
	}
	if len(rubberDucks) > 0 {
		for index, rubberDuck := range rubberDucks {
			d := distanceBetweenPoints(l.x, l.y, rubberDuck.x, rubberDuck.y)
			if d < minDist {
				minDist = d
				l.homingTargetIndex = index
				l.homingTargetType = "RubberDuck"
			}
		}
	}
}

func (l *Laser) draw(screen *ebiten.Image, x float64, y float64) {
	// Draw the laser line
	laserX := l.x - x
	laserY := l.y - y
	x1, y1 := laserX, laserY
	x2, y2 := laserX+l.size*math.Cos(l.angle), laserY+l.size*math.Sin(l.angle)
	vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 3, l.color, false)
}
