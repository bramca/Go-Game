package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TempReward struct {
	duration   int
	reward     string
	properties map[string]any
}

func closestLootBoxIndex() int {
	minDistance := distanceBetweenPoints(player.x, player.y, lootBoxes[0].x, lootBoxes[0].y)
	minIndex := 0
	for index, lootBox := range lootBoxes {
		currDistance := distanceBetweenPoints(player.x, player.y, lootBox.x, lootBox.y)
		if currDistance < minDistance {
			minIndex = index
			minDistance = currDistance
		}
	}
	return minIndex
}

func (t *TempReward) apply() {
	switch t.reward {
	case "Detect Boxes":
		t.properties["lootBoxIndex"] = -1
		t.properties["color"] = color.RGBA{255, 240, 0, 255}
	case "Invincible":
		player.healthBar.healthBarColor = color.RGBA{254, 241, 96, 255}
		player.invincible = true
	case "Insta Kill":
		t.properties["originalBackgroundColor"] = backgroundColor
		t.properties["wasVampire"] = false
		backgroundColor = color.RGBA{0, 0, 0, 255}
		player.img = playerSkullImage
		if player.vampire {
			player.img = playerVampireSkullImage
			t.properties["wasVampire"] = true
		}
		player.w = float64(player.img.Bounds().Dx())
		player.h = float64(player.img.Bounds().Dy())
		camX = player.x + player.w/2 - screenWidth/2
		camY = player.y + player.h/2 - screenHeight/2
		player.instaKill = true
	case "Vampire Mode":
		t.properties["wasInstaKill"] = false
		t.properties["originalDamageColor"] = damageColor
		t.properties["originalLaserColor"] = playerLaserColor
		playerLaserColor = color.RGBA{251, 166, 157, 255}
		damageColor = color.RGBA{255, 0, 0, 255}
		player.img = playerVampireImage
		if player.instaKill {
			player.img = playerVampireSkullImage
			t.properties["wasInstaKill"] = true
		}
		player.w = float64(player.img.Bounds().Dx())
		player.h = float64(player.img.Bounds().Dy())
		camX = player.x + player.w/2 - screenWidth/2
		camY = player.y + player.h/2 - screenHeight/2
		player.vampire = true
	}
}

func (t *TempReward) update() {
	switch t.reward {
	case "Detect Boxes":
		t.properties["lootBoxIndex"] = -1
		if len(lootBoxes) > 0 {
			t.properties["lootBoxIndex"] = closestLootBoxIndex()
		}
		t.properties["color"] = color.RGBA{255, 240, 0, 255}
		t.duration -= 1
	case "Invincible":
		t.duration -= 1
		if t.duration <= 0 {
			player.invincible = false
			player.healthBar.healthBarColor = playerHealthbarColors[0]
		}
	case "Insta Kill":
		t.duration -= 1
		if t.properties["wasVampire"].(bool) && !player.vampire {
			player.img = playerSkullImage
			player.w = float64(player.img.Bounds().Dx())
			player.h = float64(player.img.Bounds().Dy())
		}
		if t.duration <= 0 {
			player.instaKill = false
			player.img = playerImage
			if player.vampire {
				player.img = playerVampireImage
			}
			player.w = float64(player.img.Bounds().Dx())
			player.h = float64(player.img.Bounds().Dy())
			backgroundColor = t.properties["originalBackgroundColor"].(color.RGBA)
		}
	case "Vampire Mode":
		t.duration -= 1
		if t.properties["wasInstaKill"].(bool) && !player.instaKill {
			player.img = playerVampireImage
			player.w = float64(player.img.Bounds().Dx())
			player.h = float64(player.img.Bounds().Dy())
		}
		if t.duration <= 0 {
			player.vampire = false
			player.img = playerImage
			if player.instaKill {
				player.img = playerSkullImage
			}
			player.w = float64(player.img.Bounds().Dx())
			player.h = float64(player.img.Bounds().Dy())
			damageColor = t.properties["originalDamageColor"].(color.RGBA)
			playerLaserColor = t.properties["originalLaserColor"].(color.RGBA)
		}
	}
}

func (t *TempReward) draw(screen *ebiten.Image) {
	switch t.reward {
	case "Detect Boxes":
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
				vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 1, t.properties["color"].(color.Color), false)
				angle2 := 2*math.Pi - (math.Pi - angle) - math.Pi/4
				x3, y3 := x2+smallEnd*math.Cos(angle2), y2+smallEnd*math.Sin(angle2)
				vector.StrokeLine(screen, float32(x2), float32(y2), float32(x3), float32(y3), 1, t.properties["color"].(color.Color), false)
				angle3 := 2*math.Pi - (math.Pi - angle) + math.Pi/4
				x4, y4 := x2+smallEnd*math.Cos(angle3), y2+smallEnd*math.Sin(angle3)
				vector.StrokeLine(screen, float32(x2), float32(y2), float32(x4), float32(y4), 1, t.properties["color"].(color.Color), false)
			}
		}
	}
}
