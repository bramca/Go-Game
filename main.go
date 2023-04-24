package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
	ModePause
)

const (
	screenWidth  = 1280
	screenHeight = 860
)

var (
	playerStartSpeed        = 6.0
	playerStartAcceleration = 0.2
	pointsPerHit            = 2

	backgroundColor = color.RGBA{R: 8, G: 14, B: 44, A: 1}

	camX = 0.0
	camY = 0.0

	healthBarSize = 5.0

	playerStartPoints     = 15
	playerFriction        = 0.05
	playerLaserColor      = color.RGBA{R: 183, G: 244, B: 216, A: 255}
	scoreColor            = color.RGBA{R: 255, G: 255, B: 255, A: 240}
	playerStartFireRate   = framesPerSecond / 3
	playerFireFrameCount  = -1
	playerHealthbarColors = []color.RGBA{{0, 255, 0, 240}, {255, 0, 0, 240}}
	player                = &Player{
		x:            0,
		y:            0,
		w:            20,
		h:            30,
		angle:        0.0,
		points:       playerStartPoints,
		maxPoints:    playerStartPoints,
		fireRate:     playerStartFireRate,
		laserSpeed:   laserSpeed,
		speed:        playerStartSpeed,
		acceleration: playerStartAcceleration,
		damage:       pointsPerHit,
	}

	enemyImages          = []*ebiten.Image{}
	enemies              = []*Enemy{}
	enemySpawnRate       = 4 * framesPerSecond
	enemyStartPoints     = 20
	enemyHitColor        = color.RGBA{R: 255, G: 240, B: 0, A: 240}
	enemyLaserColor      = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	enemyHealthbarColors = []color.RGBA{{0, 255, 0, 240}, {255, 0, 0, 240}}
	maxEnemies           = 5

	damageColor = color.RGBA{255, 255, 255, 240}

	framesPerSecond = 60
	frameCount      = 1
	maxFrameCount   = 1200

	dotFontSize    = 8
	dots           = []*Dot{}
	dotSpawnRate   = 3 * framesPerSecond
	dotSpawnCount  = 20
	dotHexSize     = 3
	minDotDuration = 8 * framesPerSecond
	maxDotDuration = 16 * framesPerSecond
	dotHitColor    = color.RGBA{R: 147, G: 250, B: 165, A: 255}
	pointsPerDot   = 1

	lootBoxImage           *ebiten.Image
	lootBoxHealthbarColors = []color.RGBA{{108, 122, 137, 1}, backgroundColor}
	maxLootBoxes           = 5
	minLootBoxDuration     = 10 * framesPerSecond
	maxLootBoxDuration     = 20 * framesPerSecond
	lootBoxHealth          = 20
	lootBoxHitColor        = color.RGBA{R: 255, G: 240, B: 0, A: 240}
	lootBoxSpawnRate       = 6 * framesPerSecond
	lootBoxes              = []*LootBox{}
	lootRewards            = []string{}
	lootScoreReward        = 200

	tempRewardDuration = 20 * framesPerSecond

	dotTextFont     font.Face
	hitTextFont     font.Face
	scoreTextFont   font.Face
	titleArcadeFont font.Face
	arcadeFont      font.Face
	healthBarFont   font.Face

	healthBarFontColor = color.RGBA{255, 255, 255, 240}

	fontSize      = 24
	titleFontSize = 36

	hitFontSize       = 10
	scoreFontSize     = 14
	healthBarFontSize = 7

	laserSpeed    = 8.0
	laserDuration = 5 * framesPerSecond
	laserSize     = 14.0

	mouseButtonClicked = false

	recticle = Recticle{
		size: 6,
	}
)

// Game implements ebiten.Game interface.
type Game struct {
	mode Mode
}

func (g *Game) initialize() {
	img, _, _ := ebitenutil.NewImageFromFile("./resources/gopher.png")
	lootBoxImage, _, _ = ebitenutil.NewImageFromFile("./resources/github.png")
	lootRewards = []string{"Health", "Firerate", "Movement", "Damage", fmt.Sprintf("%d", lootScoreReward), "Laser Speed", "Detect Boxes", "Invincible"}
	dots = []*Dot{}
	enemies = []*Enemy{}
	lootBoxes = []*LootBox{}

	player.img = img
	player.w = float64(img.Bounds().Dx())
	player.h = float64(img.Bounds().Dy())
	player.xSpeed = 0
	player.ySpeed = 0
	player.points = playerStartPoints
	player.maxPoints = playerStartPoints
	player.healthBar = HealthBar{
		x:               player.x,
		y:               player.y - player.h,
		w:               player.w,
		h:               healthBarSize,
		points:          player.points,
		maxPoints:       player.maxPoints,
		healthBarColor:  playerHealthbarColors[0],
		healthLostColor: playerHealthbarColors[1],
		textFont:        healthBarFont,
	}
	player.lasers = []*Laser{}
	player.tempRewards = []*TempReward{}
	player.score = 0
	player.fireRate = playerStartFireRate
	player.speed = playerStartSpeed
	player.acceleration = playerStartAcceleration
	player.damage = pointsPerHit

	// Calculate the position of the screen center based on the player's position
	camX = player.x + player.w/2 - screenWidth/2
	camY = player.y + player.h/2 - screenHeight/2

	spawnDots(screenWidth, screenHeight)

	// spawnEnemies()

	spawnLootBoxes()
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.mode = ModeGame
		}
	case ModeGameOver:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.initialize()
			g.mode = ModeGame
		}
	case ModePause:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.mode = ModeGame
		}
	case ModeGame:
		if player.points <= 0 {
			g.mode = ModeGameOver
			return nil
		}
		// Write your game's logical update.
		frameCount += 1

		keyPressed := false
		if math.Sqrt(math.Pow(player.xSpeed, 2)+math.Pow(player.ySpeed, 2)) < player.speed {
			if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
				player.ySpeed += player.acceleration
				keyPressed = true
			}
			if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyZ) || ebiten.IsKeyPressed(ebiten.KeyW) {
				player.ySpeed -= player.acceleration
				keyPressed = true
			}

			if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
				player.xSpeed += player.acceleration
				keyPressed = true
			}
			if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyQ) || ebiten.IsKeyPressed(ebiten.KeyA) {
				player.xSpeed -= player.acceleration
				keyPressed = true
			}

		}

		if !keyPressed {
			if player.ySpeed != 0 {
				player.ySpeed -= (player.ySpeed / math.Abs(player.ySpeed)) * playerFriction
			}
			if player.xSpeed != 0 {
				player.xSpeed -= (player.xSpeed / math.Abs(player.xSpeed)) * playerFriction
			}
		}

		// Calculate the position of the screen center based on the player's position
		camX = player.x + player.w/2 - screenWidth/2
		camY = player.y + player.h/2 - screenHeight/2

		if frameCount%maxFrameCount == 0 {
			frameCount = 1
		}

		// Generate a set of random dots if the dots slice is empty
		if frameCount%dotSpawnRate == 0 {
			spawnDots(screenWidth*2, screenHeight*2)
		}

		if frameCount%enemySpawnRate == 0 {
			spawnEnemies()
		}

		if frameCount%lootBoxSpawnRate == 0 {
			spawnLootBoxes()
		}

		// Update dots
		for _, dot := range dots {
			if !dot.eaten {
				dot.update()
			}
		}

		// Update enemies
		for _, enemy := range enemies {
			if !enemy.dead {
				enemy.brain(dots, player)
				enemy.update()
			}
			if len(enemy.lasers) > 0 {
				enemy.updateLasers()
			}
		}

		// Update Lootboxes
		for _, lootBox := range lootBoxes {
			if !lootBox.broken {
				lootBox.update()
			} else if !lootBox.rewardGiven {
				lootBox.giveReward()
			}
		}

		// Update the player rotation based on the mouse position
		player.update(float64(player.x-camX), float64(player.y-camY), dots)

		if len(player.lasers) > 0 {
			player.updateLasers()
		}

		if len(player.tempRewards) > 0 {
			player.updateTempRewards()
		}

		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// mouseButtonClicked = false
			playerFireFrameCount = -1
		}

		// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !mouseButtonClicked {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			playerFireFrameCount += 1
			if playerFireFrameCount%player.fireRate == 0 {
				player.lasers = append(player.lasers, &Laser{
					x:        player.x,
					y:        player.y,
					angle:    player.angle,
					speed:    player.laserSpeed,
					color:    playerLaserColor,
					duration: laserDuration,
					size:     laserSize,
					damage:   player.damage,
				})
				playerFireFrameCount = 0
			}
			// mouseButtonClicked = true
		}

		if ebiten.IsKeyPressed(ebiten.KeyP) {
			g.mode = ModePause
		}
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	screen.Fill(backgroundColor)
	switch g.mode {
	case ModeTitle:
		titleTexts := []string{"GO FOREVER"}
		texts := []string{"", "", "", "", "", "", "", "PRESS SPACE KEY"}

		for i, l := range titleTexts {
			x := (screenWidth - len(l)*titleFontSize) / 2
			text.Draw(screen, l, titleArcadeFont, x, (i+4)*titleFontSize, color.White)
		}

		for i, l := range texts {
			x := (screenWidth - len(l)*fontSize) / 2
			text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
		}
		for index := len(dots) - 1; index >= 0; index-- {
			dots[index].draw(screen, camX, camY)
		}
		recticle.draw(screen)
	case ModeGameOver:
		texts := []string{"", "GAME OVER!", "", "", "PRESS SPACE KEY"}
		for i, l := range texts {
			x := (screenWidth - len(l)*fontSize) / 2
			text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
		}
		player.drawScore(screen)
		for index := len(dots) - 1; index >= 0; index-- {
			if !dots[index].eaten {
				dots[index].draw(screen, camX, camY)
			}
		}
		recticle.draw(screen)
	case ModePause:
		texts := []string{"", "PAUSED", "", "", "PRESS SPACE KEY"}
		for i, l := range texts {
			x := (screenWidth - len(l)*fontSize) / 2
			text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
		}
		// Draw the dots at their current position relative to the camera
		for index := len(dots) - 1; index >= 0; index-- {
			if !dots[index].eaten && dots[index].duration > 0 {
				dots[index].draw(screen, camX, camY)
			} else if len(dots[index].hits) > 0 {
				dots[index].drawHits(screen, camX, camY)
			}
		}

		// Draw the enemies
		for index := len(enemies) - 1; index >= 0; index-- {
			if enemies[index].points > 0 && !enemies[index].dead {
				enemies[index].draw(screen, float64(enemies[index].x-camX), float64(enemies[index].y-camY), dots)
				enemies[index].drawHits(screen)
				enemies[index].drawLasers(screen, enemies[index].x-camX, enemies[index].y-camY)
			}
		}

		for index := len(lootBoxes) - 1; index >= 0; index-- {
			if lootBoxes[index].hitpoints > 0 && !lootBoxes[index].broken && lootBoxes[index].duration > 0 {
				lootBoxes[index].draw(screen, float64(lootBoxes[index].x-camX), float64(lootBoxes[index].y-camY))
				lootBoxes[index].drawHits(screen)
			}
		}

		// Draw the lasers
		player.drawLasers(screen, camX, camY)

		// Draw the player
		player.drawScore(screen)
		player.draw(screen, float64(player.x-camX), float64(player.y-camY))
		player.drawTempRewards(screen)

		// Draw recticle
		recticle.draw(screen)
	case ModeGame:
		// Translate the screen to center it on the player
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(camX), -float64(camY))

		// Draw the dots at their current position relative to the camera
		for index := len(dots) - 1; index >= 0; index-- {
			if !dots[index].eaten && dots[index].duration > 0 {
				dots[index].draw(screen, camX, camY)
			} else if len(dots[index].hits) > 0 {
				dots[index].drawHits(screen, camX, camY)
			} else {
				dots[index] = dots[len(dots)-1]
				dots = dots[:len(dots)-1]
			}
		}

		// Draw the enemies
		for index := len(enemies) - 1; index >= 0; index-- {
			if enemies[index].points > 0 && !enemies[index].dead {
				enemies[index].draw(screen, float64(enemies[index].x-camX), float64(enemies[index].y-camY), dots)
				enemies[index].drawHits(screen)
				enemies[index].drawLasers(screen, enemies[index].x-camX, enemies[index].y-camY)
			} else if !enemies[index].dead {
				enemies[index].hits = append(enemies[index].hits, Hit{
					Dot: Dot{
						x:        int(enemies[index].x),
						y:        int(enemies[index].y - enemies[index].h),
						color:    enemyHitColor,
						msg:      "+" + strconv.Itoa(enemies[index].maxPoints),
						textFont: hitTextFont,
					},
					duration: 2 * framesPerSecond / 3,
				})
				enemies[index].dead = true
				player.score += enemies[index].maxPoints
			} else if len(enemies[index].hits) > 0 || len(enemies[index].lasers) > 0 {
				enemies[index].drawHits(screen)
				enemies[index].drawLasers(screen, enemies[index].x-camX, enemies[index].y-camY)
			} else {
				enemies[index] = enemies[len(enemies)-1]
				enemies = enemies[:len(enemies)-1]
			}
		}

		// Draw the lootboxes
		for index := len(lootBoxes) - 1; index >= 0; index-- {
			if lootBoxes[index].hitpoints > 0 && !lootBoxes[index].broken && lootBoxes[index].duration > 0 {
				lootBoxes[index].draw(screen, float64(lootBoxes[index].x-camX), float64(lootBoxes[index].y-camY))
				lootBoxes[index].drawHits(screen)
			} else if !lootBoxes[index].broken && lootBoxes[index].duration > 0 {
				lootBoxes[index].hits = append(lootBoxes[index].hits, Hit{
					Dot: Dot{
						x:        int(lootBoxes[index].x),
						y:        int(lootBoxes[index].y - lootBoxes[index].h),
						color:    lootBoxHitColor,
						msg:      "+" + lootBoxes[index].reward,
						textFont: hitTextFont,
					},
					duration: framesPerSecond,
				})
				lootBoxes[index].broken = true
			} else if len(lootBoxes[index].hits) > 0 {
				lootBoxes[index].drawHits(screen)
			} else {
				lootBoxes[index] = lootBoxes[len(lootBoxes)-1]
				lootBoxes = lootBoxes[:len(lootBoxes)-1]
			}
		}

		// Draw the lasers
		player.drawLasers(screen, camX, camY)

		// Draw the player
		player.drawScore(screen)
		player.draw(screen, float64(player.x-camX), float64(player.y-camY))
		player.drawTempRewards(screen)

		// Draw recticle
		recticle.draw(screen)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// TODO: ideas
// 1. rubber duck that runs away, multiple loot box rewards received
// 2. add temporary rewards: insta kill, vampire
func main() {
	game := &Game{}
	// Sepcify the window size as you like. Here, a doulbed size is specified.
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Go Forever")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	initialize()

	game.initialize()

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
