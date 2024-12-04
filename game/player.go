package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Player struct {
	X, Y      int
	Health    int
	Attack    int
	Inventory []*Item
	Equipped  []*Item
}

func NewPlayer(x, y int) *Player {
	return &Player{
		X:         x,
		Y:         y,
		Health:    100,
		Attack:    10,
		Inventory: []*Item{},
		Equipped:  []*Item{},
	}
}

func (p *Player) Update(g *Game) {
	dx, dy := 0, 0
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		dy = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		dy = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		dx = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		dx = 1
	}

	newX, newY := p.X+dx, p.Y+dy
	if newX >= 0 && newX < g.Map.Width && newY >= 0 && newY < g.Map.Height {
		if g.Map.Tiles[newY][newX].Walkable {
			p.X, p.Y = newX, newY
		}
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	posX := float64(p.X * TileSize)
	posY := float64(p.Y * TileSize)
	text.Draw(screen, "@", basicfont.Face7x13, int(posX), int(posY)+13, ColorWhite)
}
