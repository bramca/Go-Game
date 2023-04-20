package main

import (
	"encoding/hex"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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

func spawnDots(xBound, yBound int) {
	for i := 0; i < dotSpawnCount; i++ {
		x := int(camX + float64(rand.Intn(xBound)))
		y := int(camY + float64(rand.Intn(yBound)))
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
			duration: minDotDuration + rand.Intn(maxDotDuration-minDotDuration),
		})
	}
}

func spawnEnemies() {
	for i := len(enemies); i < maxEnemies+int(player.score/100); i++ {
		enemyImg := enemyImages[rand.Intn(len(enemyImages))]
		x := camX + float64(rand.Intn(screenWidth*2))
		y := camY + float64(rand.Intn(screenHeight*2))
		w := float64(enemyImg.Bounds().Dx())
		h := float64(enemyImg.Bounds().Dy())
		points := enemyStartPoints + player.score/100
		maxPoints := enemyStartPoints + player.score/100
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
					x:               x,
					y:               y - h,
					w:               w,
					h:               healthBarSize,
					points:          points,
					maxPoints:       maxPoints,
					healthBarColor:  enemyHealthbarColors[0],
					healthLostColor: enemyHealthbarColors[1],
					textFont:        healthBarFont,
				},
				damage: pointsPerHit + maxPoints/10,
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

func spawnLootBoxes() {
	for i := len(lootBoxes); i < maxLootBoxes; i++ {
		x := camX + float64(rand.Intn(screenWidth*4))
		y := camY + float64(rand.Intn(screenHeight*4))
		w := float64(lootBoxImage.Bounds().Dx())
		h := float64(lootBoxImage.Bounds().Dy())
		hitPoints := lootBoxHealth + player.score/100
		lootBoxes = append(lootBoxes, &LootBox{
			x:            x,
			y:            y,
			w:            w,
			h:            h,
			broken:       false,
			reward:       lootRewards[rand.Intn(len(lootRewards))],
			hitpoints:    hitPoints,
			maxHitPoints: hitPoints,
			healthBar: HealthBar{
				x:               x,
				y:               y - h,
				w:               w,
				h:               healthBarSize,
				points:          hitPoints,
				maxPoints:       hitPoints,
				healthBarColor:  lootBoxHealthbarColors[0],
				healthLostColor: lootBoxHealthbarColors[1],
				textFont:        healthBarFont,
			},
			img:      lootBoxImage,
			duration: minLootBoxDuration + rand.Intn(maxLootBoxDuration-minLootBoxDuration),
		})
	}
}

func initialize() {
	enemyImageFiles := []string{"./resources/rust.png", "./resources/cpp.png", "./resources/java.png", "./resources/haskell.png", "./resources/javascript.png", "./resources/python.png", "./resources/csharp.png"}
	for _, imgFile := range enemyImageFiles {
		enemyImg, _, _ := ebitenutil.NewImageFromFile(imgFile)
		enemyImages = append(enemyImages, enemyImg)
	}

	// Generate a set of random dots if the dots slice is empty
	dpi := 72.0
	tt, _ := opentype.Parse(fonts.PressStart2P_ttf)
	dotTextFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(dotFontSize),
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	hitTextFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(hitFontSize),
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	scoreTextFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(scoreFontSize),
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	titleArcadeFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(titleFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	arcadeFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	healthBarFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(healthBarFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
