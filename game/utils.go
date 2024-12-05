package game

import (
	"fmt"
	"image/color"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 800
	TileSize     = 32
	UIHeight     = 200 // Высота области для подсказок и UI
)

var (
	ColorWhite  = color.RGBA{255, 255, 255, 255}
	ColorRed    = color.RGBA{255, 0, 0, 255}
	ColorGreen  = color.RGBA{0, 255, 0, 255}
	ColorYellow = color.RGBA{255, 255, 0, 255}
	ColorBlue   = color.RGBA{0, 0, 255, 255}
)

// Tile представляет собой клетку на карте
type Tile struct {
	Walkable bool
	Symbol   rune
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

// Константы для типов предметов
const (
	ItemTypeWeapon = "Weapon"
	ItemTypeArmor  = "Armor"
)
