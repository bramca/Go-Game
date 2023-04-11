package main

import (
	"encoding/hex"
	"image/color"
	"log"
	"math/rand"

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

	player = &Player{
		x:     0,
		y:     0,
		w:     20,
		h:     30,
		angle: 0.0,
	}

	dotSize       = 10
	dots          = []Dot{}
	dotSpawnRate  = 180
	dotSpawnCount = 20
	dotHexSize    = 3

	textFont font.Face

	laserSpeed = 8.0

	frameCount         = 0
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

	// Generate a set of random dots if the dots slice is empty
	if frameCount%dotSpawnRate == 0 {
		spawnDots()
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
		dot.draw(screen, camX, camY)
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
	img, _, _ := ebitenutil.NewImageFromFile("./gopher.png", ebiten.FilterDefault)

	player.img = img
	player.w = float64(img.Bounds().Dx())
	player.h = float64(img.Bounds().Dy())

	// Calculate the position of the screen center based on the player's position
	camX = player.x + player.w/2 - screenWidth/2
	camY = player.y + player.h/2 - screenHeight/2

	// Generate a set of random dots if the dots slice is empty
	dpi := 72.0
	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	textFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(dotSize),
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	spawnDots()

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
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
		msg, _ := randomHex(4)
		dots = append(dots, Dot{
			x: x,
			y: y,
			color: color.RGBA{
				R: 0x80 + uint8(rand.Intn(0x7f)),
				G: 0x80 + uint8(rand.Intn(0x7f)),
				B: 0x80 + uint8(rand.Intn(0x7f)),
				A: 0xf0,
			},
			msg:      msg,
			textFont: textFont,
		})
	}
}
