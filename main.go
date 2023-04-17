package main

import (
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
)

const (
	screenWidth  = 1280
	screenHeight = 860
)

var (
	maxSpeed    = 6.0
	speedUpdate = 0.2

	camX = 0.0
	camY = 0.0

	healthBarSize = 5.0

	playerStartPoints = 15
	playerFriction    = 0.05
	playerLaserColor  = color.RGBA{R: 183, G: 244, B: 216, A: 255}
	scoreColor        = color.RGBA{R: 255, G: 255, B: 255, A: 240}
	player            = &Player{
		x:         0,
		y:         0,
		w:         20,
		h:         30,
		angle:     0.0,
		points:    playerStartPoints,
		maxPoints: playerStartPoints,
	}

	enemyImages      = []string{}
	enemies          = []*Enemy{}
	enemySpawnRate   = 4 * framesPerSecond
	enemyStartPoints = 20
	enemyHitColor    = color.RGBA{R: 255, G: 240, B: 0, A: 240}
	enemyLaserColor  = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	maxEnemies       = 5

	framesPerSecond = 60
	frameCount      = 1
	maxFrameCount   = 1200

	dotFontSize   = 8
	dots          = []*Dot{}
	dotSpawnRate  = 3 * framesPerSecond
	dotSpawnCount = 20
	dotHexSize    = 3
	dotHitColor   = color.RGBA{R: 147, G: 250, B: 165, A: 255}
	pointsPerDot  = 1

	dotTextFont     font.Face
	hitTextFont     font.Face
	scoreTextFont   font.Face
	titleArcadeFont font.Face
	arcadeFont      font.Face

	fontSize      = 24
	titleFontSize = 36

	hitFontSize   = 10
	scoreFontSize = 14

	laserSpeed    = 8.0
	laserDuration = 5 * framesPerSecond
	pointsPerHit  = 2

	mouseButtonClicked = false

	recticle = Recticle{
		size: 6,
	}

	backgroundColor = color.RGBA{R: 8, G: 14, B: 44, A: 1}
)

// Game implements ebiten.Game interface.
type Game struct {
	mode Mode
}

func (g *Game) initialize() {
	img, _, _ := ebitenutil.NewImageFromFile("./resources/gopher.png")
	dots = []*Dot{}
	enemies = []*Enemy{}

	player.img = img
	player.w = float64(img.Bounds().Dx())
	player.h = float64(img.Bounds().Dy())
	player.points = playerStartPoints
	player.maxPoints = playerStartPoints
	player.healthBar = HealthBar{
		x:         player.x,
		y:         player.y - player.h,
		w:         player.w,
		h:         healthBarSize,
		points:    player.points,
		maxPoints: player.maxPoints,
	}
	player.score = 0

	// Calculate the position of the screen center based on the player's position
	camX = player.x + player.w/2 - screenWidth/2
	camY = player.y + player.h/2 - screenHeight/2

	spawnDots(screenWidth, screenHeight)

	spawnEnemies()
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
	case ModeGame:
		if player.points <= 0 {
			g.mode = ModeGameOver
			return nil
		}
		// Write your game's logical update.
		frameCount += 1

		keyPressed := false
		if math.Sqrt(math.Pow(player.xSpeed, 2)+math.Pow(player.ySpeed, 2)) < maxSpeed {
			if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
				player.ySpeed += speedUpdate
				keyPressed = true
			}
			if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyZ) || ebiten.IsKeyPressed(ebiten.KeyW) {
				player.ySpeed -= speedUpdate
				keyPressed = true
			}

			if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
				player.xSpeed += speedUpdate
				keyPressed = true
			}
			if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyQ) || ebiten.IsKeyPressed(ebiten.KeyA) {
				player.xSpeed -= speedUpdate
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

		// Update the player rotation based on the mouse position
		player.update(float64(player.x-camX), float64(player.y-camY), dots)

		if len(player.lasers) > 0 {
			player.updateLasers()
		}

		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mouseButtonClicked = false
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !mouseButtonClicked {
			player.lasers = append(player.lasers, &Laser{
				x:        player.x,
				y:        player.y,
				angle:    player.angle,
				speed:    laserSpeed,
				color:    playerLaserColor,
				duration: laserDuration,
			})
			mouseButtonClicked = true
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
	case ModeGame:
		// Translate the screen to center it on the player
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(camX), -float64(camY))

		// Draw the dots at their current position relative to the camera
		for index := len(dots) - 1; index >= 0; index-- {
			if !dots[index].eaten {
				dots[index].draw(screen, camX, camY)
			} else if len(dots[index].hits) > 0 {
				dots[index].drawHits(screen, camX, camY)
			} else {
				dots[index] = dots[len(dots)-1]
				dots = dots[:len(dots)-1]
			}
		}

		// Draw the enemies
		for index, enemy := range enemies {
			if enemy.points > 0 && !enemy.dead {
				enemy.draw(screen, float64(enemy.x-camX), float64(enemy.y-camY), dots)
				enemy.drawHits(screen)
				enemy.drawLasers(screen, enemy.x-camX, enemy.y-camY)
			} else if !enemy.dead {
				enemy.hits = append(enemy.hits, Hit{
					Dot: Dot{
						x:        int(enemy.x),
						y:        int(enemy.y - enemy.h),
						color:    enemyHitColor,
						msg:      "+" + strconv.Itoa(enemy.maxPoints),
						textFont: hitTextFont,
					},
					duration: 2 * framesPerSecond / 3,
				})
				enemy.dead = true
				player.score += enemy.maxPoints
			} else if len(enemy.hits) > 0 || len(enemy.lasers) > 0 {
				enemy.drawHits(screen)
				enemy.drawLasers(screen, enemy.x-camX, enemy.y-camY)
			} else {
				enemies[index] = enemies[len(enemies)-1]
				enemies = enemies[:len(enemies)-1]
			}
		}

		// Draw the lasers
		player.drawLasers(screen, camX, camY)

		// Draw the player
		player.drawScore(screen)
		player.draw(screen, float64(player.x-camX), float64(player.y-camY))

		// Draw recticle
		recticle.draw(screen)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

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
