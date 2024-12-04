package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player *Player
	Map    *GameMap
}

func NewGame() *Game {
	gm := NewGameMap()
	player := NewPlayer(gm.PlayerStartX, gm.PlayerStartY)
	return &Game{
		Player: player,
		Map:    gm,
	}
}

func (g *Game) Update() error {
	g.Player.Update(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Map.Draw(screen)
	g.Player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
