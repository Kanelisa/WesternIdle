package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"WesternIdle/internal/core"
)

func main() {
	game := core.NewGameInstance()

	ebiten.SetWindowSize(core.ScreenWidth, core.ScreenHeight)
	ebiten.SetWindowTitle("Western Idle")
	ebiten.SetTPS(60)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
