package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 1280
	screenHeight = 860
)

var (
	maxSpeed = 5.0

	camX = 0.0
	camY = 0.0

	player = &Player{
		x:     0,
		y:     0,
		w:     20,
		h:     30,
		angle: 0.0,
	}

	dotSize      = 4
	dots         = []DrawRectParams{}
	dotSpawnRate = 120
	dotSpawnCount = 30

	laserSpeed = 8.0

	frameCount = 0
	mouseButtonClicked = false

	recticle = Recticle{
		size: 5,
	}
)

type DrawRectParams struct {
	X     float64
	Y     float64
	W     float64
	H     float64
	Color color.RGBA
}

// DrawLineParams{X0: laserX - float64(camX), Y0: laserY - float64(camY), X1: laserEndX - float64(camX), Y1: laserEndY - float64(camY), Color: color.RGBA{255, 0, 0, 255}})
type DrawLineParams struct {
	X0    float64
	X1    float64
	Y0    float64
	Y1    float64
	Color color.RGBA
}

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	// Write your game's logical update.
	frameCount += 1
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if player.speed < maxSpeed {
			player.speed += 0.1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if player.speed > -maxSpeed {
			player.speed -= 0.1
		}
	}

	// Calculate the position of the screen center based on the player's position
	camX = player.x + player.w/2 - screenWidth/2
	camY = player.y + player.h/2 - screenHeight/2

	// Generate a set of random dots if the dots slice is empty
	if frameCount%dotSpawnRate == 0 {
		for i := 0; i < dotSpawnCount; i++ {
			x := camX + float64(rand.Intn(screenWidth))
			y := camY + float64(rand.Intn(screenHeight))
			dots = append(dots, DrawRectParams{X: float64(x), Y: float64(y), W: float64(dotSize), H: float64(dotSize), Color: color.RGBA{0, 255, 0, 255}})
		}
		frameCount = 1
	}

	// Update the player rotation based on the mouse position
	player.update(float64(player.x-camX), float64(player.y-camY))

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
	for _, dot := range dots {
		ebitenutil.DrawRect(screen, dot.X-float64(camX), dot.Y-float64(camY), dot.W, dot.H, dot.Color)
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
	ebiten.SetWindowTitle("Your game's title")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	img, _, _ := ebitenutil.NewImageFromFile("./spaceship.gif", ebiten.FilterDefault)

	player.img = img
	player.w = float64(img.Bounds().Dx())
	player.h = float64(img.Bounds().Dy())

	// Calculate the position of the screen center based on the player's position
	camX = player.x + player.w/2 - screenWidth/2
	camY = player.y + player.h/2 - screenHeight/2

	// Generate a set of random dots if the dots slice is empty
	for i := 0; i < dotSpawnCount; i++ {
		x := camX + float64(rand.Intn(screenWidth))
		y := camY + float64(rand.Intn(screenHeight))
		dots = append(dots, DrawRectParams{X: float64(x), Y: float64(y), W: float64(dotSize), H: float64(dotSize), Color: color.RGBA{0, 255, 0, 255}})
	}

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
