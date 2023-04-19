package main

import "github.com/hajimehoshi/ebiten/v2"

type LootBox struct {
	x, y, w, h float64
	hits       []Hit
	broken     bool
	reward     string
	hitpoints  int
	healthBar  HealthBar
	img        *ebiten.Image
}

// TODO: do the effect of the reward on broken
func (l *LootBox) update() {
	l.healthBar.update(l.x-camX-l.w/2, l.y-(l.h-l.h/3)-camY, l.hitpoints, lootBoxHealth)
}

func (l *LootBox) drawHits(screen *ebiten.Image) {
	for i := len(l.hits) - 1; i >= 0; i-- {
		if l.hits[i].duration > 0 {
			l.hits[i].update()
			l.hits[i].draw(screen, camX, camY)
		} else {
			l.hits[i] = l.hits[len(l.hits)-1]
			l.hits = l.hits[:len(l.hits)-1]
		}
	}
}

func (l *LootBox) draw(screen *ebiten.Image, x float64, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(l.w/2), -float64(l.h/2))
	op.GeoM.Translate(x, y)
	screen.DrawImage(l.img, op)
	l.healthBar.draw(screen)
}
