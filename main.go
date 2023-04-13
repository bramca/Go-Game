package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 1280
	screenHeight = 860
)

var (
	maxSpeed    = 5.0
	speedUpdate = 0.2

	camX = 0.0
	camY = 0.0

	healthBarSize = 5.0

	playerStartPoints = 15
	player            = &Player{
		x:         0,
		y:         0,
		w:         20,
		h:         30,
		angle:     0.0,
		points:    playerStartPoints,
		maxPoints: playerStartPoints,
	}

	enemyImages = []string{}
	enemies          = []*Enemy{}
	enemySpawnRate = 4 * framesPerSecond
	enemyStartPoints = 20
	maxEnemies       = 5

	framesPerSecond = 60
	frameCount      = 1
	maxFrameCount   = 1200

	dotSize       = 10
	dots          = []*Dot{}
	dotSpawnRate  = 3 * framesPerSecond
	dotSpawnCount = 20
	dotHexSize    = 3
	pointsPerDot  = 1

	dotTextFont font.Face
	hitTextFont font.Face

	hitSize = 13

	laserSpeed   = 8.0
	maxLasers    = 10
	pointsPerHit = 2

	mouseButtonClicked = false

	recticle = Recticle{
		size: 6,
	}
)

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	// Write your game's logical update.
	frameCount += 1
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if player.ySpeed < maxSpeed {
			player.ySpeed += speedUpdate
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if player.ySpeed > -maxSpeed {
			player.ySpeed -= speedUpdate
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if player.xSpeed < maxSpeed {
			player.xSpeed += speedUpdate
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if player.xSpeed > -maxSpeed {
			player.xSpeed -= speedUpdate
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
		spawnDots()
	}

	if frameCount%enemySpawnRate == 0 {
		spawnEnemies()
	}

	// Update enemies
	for _, enemy := range enemies {
		enemy.brain(dots, player)
		enemy.update()
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
			x:     player.x,
			y:     player.y,
			angle: player.angle,
			speed: laserSpeed,
		})
		if len(player.lasers) > maxLasers {
			player.lasers[0] = player.lasers[len(player.lasers)-1]
			player.lasers = player.lasers[1:]
		}
		mouseButtonClicked = true
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	screen.Fill(color.Black)

	// Translate the screen to center it on the player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(camX), -float64(camY))

	// Draw the dots at their current position relative to the camera
	for index := len(dots) - 1; index >= 0; index-- {
		if dots[index] != nil {
			dots[index].draw(screen, camX, camY)
		} else {
			dots[index] = dots[len(dots)-1]
			dots = dots[:len(dots)-1]
		}
	}

	// Draw the enemies
	for index, enemy := range enemies {
		if enemy.points > 0 {
			enemy.draw(screen, float64(enemy.x-camX), float64(enemy.y-camY), dots)
			enemy.drawLasers(screen, enemy.x-camX, enemy.y-camY)
		} else {
			enemies[index] = enemies[len(enemies)-1]
			enemies = enemies[:len(enemies)-1]
		}
	}

	// Draw the lasers
	player.drawLasers(screen, camX, camY)

	// Draw the player
	player.draw(screen, float64(player.x-camX), float64(player.y-camY))

	// Draw recticle
	recticle.draw(screen)
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
	img, _, _ := ebitenutil.NewImageFromFile("./resources/gopher.png", ebiten.FilterDefault)

	player.img = img
	player.w = float64(img.Bounds().Dx())
	player.h = float64(img.Bounds().Dy())
	player.healthBar = HealthBar{
		x:         player.x,
		y:         player.y - player.h,
		w:         player.w,
		h:         healthBarSize,
		points:    player.points,
		maxPoints: player.maxPoints,
	}

	// Calculate the position of the screen center based on the player's position
	camX = player.x + player.w/2 - screenWidth/2
	camY = player.y + player.h/2 - screenHeight/2

	enemyImages = []string{"./resources/rust.png", "./resources/cpp.png", "./resources/java.png", "./resources/haskell.png", "./resources/javascript.png", "./resources/python.png", "./resources/csharp.png"}

	// Generate a set of random dots if the dots slice is empty
	dpi := 72.0
	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	dotTextFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(dotSize),
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	hitTextFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(hitSize),
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})

	spawnDots()

	spawnEnemies()

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
