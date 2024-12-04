package game

import (
	"image/color"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	TileSize     = 16
)

var (
	ColorWhite = color.RGBA{255, 255, 255, 255}
)

type Tile struct {
	Walkable bool
	Symbol   rune
}
