package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)


type Player struct {
    x, y  float64
    w, h  float64
    angle float64
}

func (p *Player) update(x, y float64) {
    // Move the player based on the mouse position
    mx, my := ebiten.CursorPosition()
    p.angle = angleBetweenPoints(x, y, float64(mx), float64(my))
}

func (p *Player) draw(screen *ebiten.Image, x float64, y float64) {
    // Draw the player
    img, _, _ := ebitenutil.NewImageFromFile("./spaceship.gif", ebiten.FilterDefault)
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(-float64(img.Bounds().Dx()/2), -float64(img.Bounds().Dy()/2))
    op.GeoM.Rotate(p.angle + math.Pi/2)
    op.GeoM.Translate(x, y)
    screen.DrawImage(img, op)
}
func angleBetweenPoints(x1, y1, x2, y2 float64) float64 {
    return math.Atan2(y2-y1, x2-x1)
}

func rotatePoint(x, y, cx, cy, angle float64) (float64, float64) {
    s := math.Sin(angle)
    c := math.Cos(angle)

    // Translate to origin
    x -= cx
    y -= cy

    // Rotate
    x, y = x*c-y*s, x*s+y*c

    // Translate back
    x += cx
    y += cy

    return x, y
}
