package main

import (
	"log"

	"github.com/bramca/gogame"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: ideas
// 1. add more temporary rewards
// 2. add more interesting enemies
func main() {
	game := &gogame.Game{}
	// Sepcify the window size as you like. Here, a doulbed size is specified.
	ebiten.SetWindowSize(gogame.ScreenWidth, gogame.ScreenHeight)
	ebiten.SetWindowTitle("Go Forever")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	gogame.Initialize()

	game.Initialize()

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
