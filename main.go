package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	playerSp = 4.0

	player = &Player{
		x:     0,
		y:     0,
		w:     20,
		h:     30,
		angle: 0.0,
	}

	dotSize  = 4
	dots     = []DrawRectParams{}
	dotSpawnRate = 120

	laserSpeed = 8
	lasers     = []DrawLineParams{}

	mouseX = 0
	mouseY = 0

	frameCount = 1
)

type Laser struct {
	x     float64
	y     float64
	angle float64
}

type DrawRectParams struct {
	X float64
	Y float64
	W float64
	H float64
	Color color.RGBA
}
// DrawLineParams{X0: laserX - float64(camX), Y0: laserY - float64(camY), X1: laserEndX - float64(camX), Y1: laserEndY - float64(camY), Color: color.RGBA{255, 0, 0, 255}})
type DrawLineParams struct {
	X0 float64
	X1 float64
	Y0 float64
	Y1 float64
	Color color.RGBA
}

func update(screen *ebiten.Image) error {
	frameCount += 1
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		player.y -= playerSp
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		player.y += playerSp
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		player.x -= playerSp
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		player.x += playerSp
	}

	// Calculate the position of the screen center based on the player's position
	camX := player.x + player.w/2 - screenWidth/2
	camY := player.y + player.h/2 - screenHeight/2

	// Generate a set of random dots if the dots slice is empty
	if frameCount % dotSpawnRate == 0 {
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 10; i++ {
			x := camX + float64(rand.Intn(screenWidth))
			y := camY + float64(rand.Intn(screenHeight))
			dots = append(dots, DrawRectParams{X: float64(x), Y: float64(y), W: float64(dotSize), H: float64(dotSize), Color: color.RGBA{0, 255, 0, 255}})
		}
		frameCount = 1
	}

	// Update the player rotation based on the mouse position
	player.update(float64(player.x-camX), float64(player.y-camY))

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		laserX := float64(player.x) + float64(player.w)/2
		laserY := float64(player.y) + float64(player.h)/2
		laserLen := 50.0
		laserEndX := laserX + laserLen*math.Cos(player.angle)
		laserEndY := laserY + laserLen*math.Sin(player.angle)
		lasers = append(lasers, DrawLineParams{X0: laserX - float64(camX), Y0: laserY - float64(camY), X1: laserEndX - float64(camX), Y1: laserEndY - float64(camY), Color: color.RGBA{255, 0, 0, 255}})
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}


	// Clear the screen to white
	screen.Fill(color.Black)

	// Translate the screen to center it on the player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(camX), -float64(camY))

	// Draw the dots at their current position relative to the camera
	for _, dot := range dots {
		ebitenutil.DrawRect(screen, dot.X-float64(camX), dot.Y-float64(camY), dot.W, dot.H, dot.Color)
	}

	// Move the lasers
	for i := 0; i < len(lasers); i++ {
		ebitenutil.DrawLine(screen, lasers[i].X0, lasers[i].Y0, lasers[i].X1, lasers[i].Y1, lasers[i].Color)
	}

	player.draw(screen, float64(player.x-camX), float64(player.y-camY))

    // Move the player based on the mouse position
    mx, my := ebiten.CursorPosition()
    angle := angleBetweenPoints(player.x, player.y, float64(mx), float64(my))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("angle: %f", angle))

	return nil
}

func main() {
	// Set up the game window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("My Game")

	// Start the game loop
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "My Game"); err != nil {
		log.Fatal(err)
	}
}
