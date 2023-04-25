package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type LootBox struct {
	x, y, w, h   float64
	hits         []Hit
	broken       bool
	rewardGiven  bool
	reward       string
	hitpoints    int
	maxHitPoints int
	healthBar    HealthBar
	img          *ebiten.Image
	duration     int
}

func (l *LootBox) giveReward() {
	switch l.reward {
	// Health Boost
	case lootRewards[0]:
		player.points += player.maxPoints / 3
		if player.points > player.maxPoints {
			player.maxPoints = player.points
		}
		// Firerate Increase
	case lootRewards[1]:
		if player.fireRate*3/4 > 0 {
			player.fireRate = player.fireRate * 4 / 5
		}
		// Movement Increase
	case lootRewards[2]:
		player.speed += 0.2
		player.acceleration += 0.02
		// Damage Increase
	case lootRewards[3]:
		player.damage += pointsPerHit
		// Score Increase
	case lootRewards[4]:
		player.score += lootScoreReward
		// Laser Speed
	case lootRewards[5]:
		player.laserSpeed += 2.0
		// Detect Boxes
	case lootRewards[6]:
		activateTempReward(l.reward, tempRewardDuration)
		// Invincible
	case lootRewards[7]:
		activateTempReward(l.reward, tempRewardDuration)
		// Insta Kill
	case lootRewards[8]:
		activateTempReward(l.reward, tempRewardDuration)
		// Vampire Mode
	case lootRewards[9]:
		activateTempReward(l.reward, tempRewardDuration)
	}
	l.rewardGiven = true
}

func (l *LootBox) update() {
	l.duration -= 1
	l.healthBar.update(l.x-camX-l.w/2, l.y-(l.h-l.h/3)-camY, l.hitpoints, l.maxHitPoints)
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
