package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TempReward struct {
	duration   int
	reward     string
	properties map[string]any
}

func (t *TempReward) update() {
	switch t.reward {
	// Detect Boxes
	case lootRewards[6]:
		lootBoxIndex, ok := -1, false
		if lootBoxIndex, ok = t.properties["lootBoxIndex"].(int); !ok {
			t.properties["lootBoxIndex"] = rand.Intn(len(lootBoxes))
		}
		if (lootBoxIndex > len(lootBoxes)-1 || lootBoxIndex < 0) && len(lootBoxes) > 0 {
			t.properties["lootBoxIndex"] = rand.Intn(len(lootBoxes))
		}
		t.properties["color"] = color.RGBA{255, 240, 0, 255}
		t.duration -= 1
		// Invincible
	case lootRewards[7]:
		player.healthBar.healthBarColor = color.RGBA{254, 241, 96, 255}
		t.duration -= 1
		player.invincible = true
		if t.duration <= 0 {
			player.invincible = false
			player.healthBar.healthBarColor = playerHealthbarColors[0]
		}
	}
}

func (t *TempReward) draw(screen *ebiten.Image) {
	switch t.reward {
	// Detect Boxes
	case lootRewards[6]:
		if lootBoxIndex, ok := t.properties["lootBoxIndex"].(int); ok {
			if lootBoxIndex <= len(lootBoxes)-1 && lootBoxIndex >= 0 {
				longEnd := 18.0
				smallEnd := 8.0
				lootBox := lootBoxes[lootBoxIndex]
				angle := angleBetweenPoints(player.x, player.y, lootBox.x, lootBox.y)
				startX := player.x - camX
				startY := player.y - camY
				x1, y1 := startX, startY-player.h
				x2, y2 := x1+longEnd*math.Cos(angle), y1+longEnd*math.Sin(angle)
				ebitenutil.DrawLine(screen, x1, y1, x2, y2, t.properties["color"].(color.Color))
				angle2 := 2*math.Pi - (math.Pi - angle) - math.Pi/4
				x3, y3 := x2+smallEnd*math.Cos(angle2), y2+smallEnd*math.Sin(angle2)
				ebitenutil.DrawLine(screen, x2, y2, x3, y3, t.properties["color"].(color.Color))
				angle3 := 2*math.Pi - (math.Pi - angle) + math.Pi/4
				x4, y4 := x2+smallEnd*math.Cos(angle3), y2+smallEnd*math.Sin(angle3)
				ebitenutil.DrawLine(screen, x2, y2, x4, y4, t.properties["color"].(color.Color))
			}
		}
	case lootRewards[7]:
		ebitenutil.DebugPrint(screen, fmt.Sprintf("duration: %d", t.duration))
	}
}