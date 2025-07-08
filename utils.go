package main

import (
	randcrypto "crypto/rand"
	"encoding/hex"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func angleBetweenPoints(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y2-y1, x2-x1)
}

func distanceBetweenPoints(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := randcrypto.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func setDotDrawOptions(d *Dot) {
	drawOptions := text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{},
	}
	drawOptions.ColorScale.SetR(float32(d.color.R) / 256.0)
	drawOptions.ColorScale.SetG(float32(d.color.G) / 256.0)
	drawOptions.ColorScale.SetB(float32(d.color.B) / 256.0)
	drawOptions.ColorScale.SetA(float32(d.color.G) / 256.0)
	d.drawOptions = &drawOptions
}

func spawnDots(xBound, yBound int) {
	for range dotSpawnCount {
		x := int(camX + float64(rand.Intn(xBound)))
		y := int(camY + float64(rand.Intn(yBound)))
		msg, _ := randomHex(4)
		dot := &Dot{
			x: x,
			y: y,
			color: color.RGBA{
				R: 0x80 + uint8(rand.Intn(0x7f)),
				G: 0x80 + uint8(rand.Intn(0x7f)),
				B: 0x80 + uint8(rand.Intn(0x7f)),
				A: 0xf0,
			},
			msg:      msg,
			textFont: text.NewGoXFace(dotTextFont),
			duration: minDotDuration + rand.Intn(maxDotDuration-minDotDuration),
		}
		setDotDrawOptions(dot)
		dots = append(dots, dot)
	}
}

func spawnEnemies() {
	for i := len(enemies); i < maxEnemies+int(player.score/100); i++ {
		numberOfEnemies := maxEnemies + int(player.score/100)
		enemyImg := enemyImages[rand.Intn(len(enemyImages))]
		x := camX + float64(rand.Intn(screenWidth*(4+numberOfEnemies/30)))
		y := camY + float64(rand.Intn(screenHeight*(4+numberOfEnemies/30)))
		w := float64(enemyImg.Bounds().Dx())
		h := float64(enemyImg.Bounds().Dy())
		points := enemyStartPoints + player.score/100
		maxPoints := enemyStartPoints + player.score/100
		visibleRange := float64(int(math.Min(screenWidth, screenHeight))+rand.Intn(int(math.Max(screenWidth, screenHeight))-int(math.Min(screenWidth, screenHeight)))) / 2
		aggressiveness := 0.6
		greediness := 0.4
		healthBar := HealthBar{
			x:               x,
			y:               y - h,
			w:               w,
			h:               healthBarSize,
			points:          points,
			maxPoints:       maxPoints,
			healthBarColor:  enemyHealthbarColors[0],
			healthLostColor: enemyHealthbarColors[1],
			textFont:        text.NewGoXFace(healthBarFont),
		}
		healthBar.setDrawOptions()
		enemies = append(enemies, &Enemy{
			Player: Player{
				x:          x,
				y:          y,
				w:          w,
				h:          h,
				angle:      0,
				lasers:     []*Laser{},
				laserSpeed: player.speed + float64(rand.Intn(int(laserSpeed))),
				img:        enemyImg,
				ySpeed:     0,
				xSpeed:     0,
				points:     points,
				maxPoints:  maxPoints,
				healthBar:  healthBar,
				damage:     pointsPerHit + maxPoints/10,
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

func spawnRubberDucks() {
	for i := len(rubberDucks); i < maxRubberDucks; i++ {
		x := camX + float64(rand.Intn(screenWidth*6))
		y := camY + float64(rand.Intn(screenHeight*6))
		w := float64(rubberDuckImage.Bounds().Dx())
		h := float64(rubberDuckImage.Bounds().Dy())
		points := rubberDuckStartPoints + player.score/100
		maxPoints := rubberDuckStartPoints + player.score/100
		visibleRange := float64(int(math.Min(screenWidth, screenHeight))+rand.Intn(int(math.Max(screenWidth, screenHeight))-int(math.Min(screenWidth, screenHeight)))) / 4
		healthBar := HealthBar{
			x:               x,
			y:               y - h,
			w:               w,
			h:               healthBarSize,
			points:          points,
			maxPoints:       maxPoints,
			healthBarColor:  rubberDuckHealthBarColors[0],
			healthLostColor: rubberDuckHealthBarColors[1],
			textFont:        text.NewGoXFace(healthBarFont),
		}
		healthBar.setDrawOptions()
		rubberDucks = append(rubberDucks, &RubberDuck{
			Player: Player{
				x:           x,
				y:           y,
				w:           w,
				h:           h,
				angle:       0,
				lasers:      []*Laser{},
				tempRewards: []*TempReward{},
				img:         rubberDuckImage,
				ySpeed:      0,
				xSpeed:      0,
				points:      points,
				maxPoints:   maxPoints,
				healthBar:   healthBar,
			},
			visibleRange:       visibleRange,
			fleeRange:          2 * visibleRange,
			dotTargetIndex:     -1,
			speedMultiplyer:    (3 + rand.Intn(6)),
			movementPrediction: float64(10 + rand.Intn(30)),
			reward:             rubberDuckRewards[rand.Intn(len(rubberDuckRewards))],
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
		healthBar := HealthBar{
			x:               x,
			y:               y - h,
			w:               w,
			h:               healthBarSize,
			points:          hitPoints,
			maxPoints:       hitPoints,
			healthBarColor:  lootBoxHealthbarColors[0],
			healthLostColor: lootBoxHealthbarColors[1],
			textFont:        text.NewGoXFace(healthBarFont),
		}
		healthBar.setDrawOptions()
		lootBoxes = append(lootBoxes, &LootBox{
			x:            x,
			y:            y,
			w:            w,
			h:            h,
			broken:       false,
			reward:       lootRewards[rand.Intn(len(lootRewards))],
			hitpoints:    hitPoints,
			maxHitPoints: hitPoints,
			healthBar:    healthBar,
			img:          lootBoxImage,
			duration:     minLootBoxDuration + rand.Intn(maxLootBoxDuration-minLootBoxDuration),
		})
	}
}

func initialize() {
	enemyImageFiles := []string{"./resources/rust.png", "./resources/cpp.png", "./resources/java.png", "./resources/haskell.png", "./resources/javascript.png", "./resources/python.png", "./resources/csharp.png"}
	for _, imgFile := range enemyImageFiles {
		enemyImg, _, _ := ebitenutil.NewImageFromFile(imgFile)
		enemyImages = append(enemyImages, enemyImg)
	}

	lootBoxImage, _, _ = ebitenutil.NewImageFromFile("./resources/github.png")

	playerImage, _, _ = ebitenutil.NewImageFromFile("./resources/gopher.png")

	playerSkullImage, _, _ = ebitenutil.NewImageFromFile("./resources/gopher_skull.png")

	playerVampireImage, _, _ = ebitenutil.NewImageFromFile("./resources/gopher_vampire.png")

	playerVampireSkullImage, _, _ = ebitenutil.NewImageFromFile("./resources/gopher_vampire_skull.png")

	rubberDuckImage, _, _ = ebitenutil.NewImageFromFile("./resources/rubber_duck.png")

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

	titleFontColorScale.ScaleWithColor(color.White)
}

func activateTempReward(lootReward string, duration int) {
	rewardActive := false
	for _, reward := range player.tempRewards {
		if reward.reward == lootReward {
			reward.duration = tempRewardDuration
			rewardActive = true
		}
	}
	if !rewardActive {
		player.tempRewards = append(player.tempRewards, &TempReward{
			duration:   duration,
			reward:     lootReward,
			properties: map[string]any{},
		})
		player.tempRewards[len(player.tempRewards)-1].apply()
	}

}
