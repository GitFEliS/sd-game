package game

import (
	"bufio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"log"
	"os"
)

type GameMap struct {
	Tiles                      [][]Tile
	Width, Height              int
	PlayerStartX, PlayerStartY int
	Items                      []*Item
	Monsters                   []*Monster
}

func NewGameMap() *GameMap {
	gm := &GameMap{}
	gm.LoadMapFromFile("assets/levels/level1.txt")
	return gm
}

func (gm *GameMap) LoadMapFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tiles [][]Tile
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		var row []Tile
		for x, char := range line {
			tile := Tile{}
			switch char {
			case '#':
				tile.Walkable = false
				tile.Symbol = '#'
			case '.':
				tile.Walkable = true
				tile.Symbol = '.'
			case '@':
				tile.Walkable = true
				tile.Symbol = '.'
				gm.PlayerStartX = x
				gm.PlayerStartY = y
			default:
				tile.Walkable = true
				tile.Symbol = '.'
			}
			row = append(row, tile)
		}
		tiles = append(tiles, row)
		y++
	}
	gm.Tiles = tiles
	gm.Width = len(tiles[0])
	gm.Height = len(tiles)
}

func (gm *GameMap) Draw(screen *ebiten.Image) {
	for y, row := range gm.Tiles {
		for x, tile := range row {
			posX := float64(x * TileSize)
			posY := float64(y * TileSize)
			text.Draw(screen, string(tile.Symbol), basicfont.Face7x13, int(posX), int(posY)+13, ColorWhite)
		}
	}
}
