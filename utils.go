package main

import (
	"encoding/hex"
	"image/color"
	"math"
	"math/rand"
)

func angleBetweenPoints(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y2-y1, x2-x1)
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
		// x := int(camX + float64(rand.Intn(screenWidth*2)))
		// y := int(camY + float64(rand.Intn(screenHeight*2)))
		x := int(camX + float64(rand.Intn(screenWidth)))
		y := int(camY + float64(rand.Intn(screenHeight)))
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
			textFont: textFont,
		})
	}
}
