package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"

	"sd-game/game"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Roguelike Adventure")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
