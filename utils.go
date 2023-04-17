package main

import (
	"encoding/hex"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func angleBetweenPoints(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y2-y1, x2-x1)
}

func distanceBetweenPoints(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func spawnDots() {
	for i := 0; i < dotSpawnCount; i++ {
		x := int(camX + float64(rand.Intn(screenWidth*2)))
		y := int(camY + float64(rand.Intn(screenHeight*2)))
		// x := int(camX + float64(rand.Intn(screenWidth)))
		// y := int(camY + float64(rand.Intn(screenHeight)))
		msg, _ := randomHex(4)
		dots = append(dots, &Dot{
			x: x,
			y: y,
			color: color.RGBA{
				R: 0x80 + uint8(rand.Intn(0x7f)),
				G: 0x80 + uint8(rand.Intn(0x7f)),
				B: 0x80 + uint8(rand.Intn(0x7f)),
				A: 0xf0,
			},
			msg:      msg,
			textFont: dotTextFont,
		})
	}
}

func spawnEnemies() {
	for i := len(enemies); i < maxEnemies+int(player.score/100); i++ {
		enemyImg, _, _ := ebitenutil.NewImageFromFile(enemyImages[rand.Intn(len(enemyImages))], ebiten.FilterDefault)
		x := camX + float64(rand.Intn(screenWidth*2))
		y := camY + float64(rand.Intn(screenHeight*2))
		w := float64(enemyImg.Bounds().Dx())
		h := float64(enemyImg.Bounds().Dy())
		points := enemyStartPoints
		maxPoints := enemyStartPoints
		visibleRange := float64(int(math.Min(screenWidth, screenHeight))+rand.Intn(int(math.Max(screenWidth, screenHeight))-int(math.Min(screenWidth, screenHeight)))) / 2
		aggressiveness := 0.6
		greediness := 0.4
		enemies = append(enemies, &Enemy{
			Player: Player{
				x:         x,
				y:         y,
				w:         w,
				h:         h,
				angle:     0,
				lasers:    []*Laser{},
				img:       enemyImg,
				ySpeed:    0,
				xSpeed:    0,
				points:    points,
				maxPoints: maxPoints,
				healthBar: HealthBar{
					x:         x,
					y:         y - h,
					w:         w,
					h:         healthBarSize,
					points:    points,
					maxPoints: maxPoints,
				},
			},
			dotTargetIndex:     -1,
			visibleRange:       visibleRange,
			shootRange:         (1 - aggressiveness) * visibleRange,
			greedy:             greediness,
			aggressive:         aggressiveness,
			shootFreq:          (1 + rand.Intn(3)) * (framesPerSecond / 4),
			speedMultiplyer:    (2 + rand.Intn(4)),
			movementPrediction: float64(10 + rand.Intn(30)),
		})
	}
}
