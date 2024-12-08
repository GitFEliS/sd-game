package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math/rand"
	"time"
)

type GameMap struct {
	Tiles                      [][]Tile
	Width, Height              int
	PlayerStartX, PlayerStartY int
	ExitX, ExitY               int
	Items                      []*Item
	Monsters                   []*Monster
}

func NewGameMap(level int) *GameMap {
	gm := &GameMap{
		Width:  25, // Размер карты по ширине
		Height: 18, // Размер карты по высоте
	}
	gm.GenerateDungeon()
	gm.SpawnMonsters(5 + level*2) // Увеличение количества монстров с каждым уровнем
	gm.SpawnItems(3 + level)      // Увеличение количества предметов с каждым уровнем
	return gm
}

func (gm *GameMap) GenerateDungeon() {
	rand.Seed(time.Now().UnixNano())

	// Инициализация всех клеток как непроходимых
	gm.Tiles = make([][]Tile, gm.Height)
	for y := 0; y < gm.Height; y++ {
		gm.Tiles[y] = make([]Tile, gm.Width)
		for x := 0; x < gm.Width; x++ {
			gm.Tiles[y][x] = Tile{
				Walkable: false,
				Symbol:   '#',
			}
		}
	}

	// Генерация комнат в последовательности
	numRooms := rand.Intn(3) + 5 // От 5 до 7 комнат
	rooms := make([]Rectangle, 0, numRooms)
	var previousRoom Rectangle

	for i := 0; i < numRooms; i++ {
		roomWidth := rand.Intn(6) + 4
		roomHeight := rand.Intn(4) + 4
		roomX := rand.Intn(gm.Width - roomWidth - 1)
		roomY := rand.Intn(gm.Height - roomHeight - 1)

		newRoom := Rectangle{X: roomX, Y: roomY, Width: roomWidth, Height: roomHeight}

		// Проверка на пересечение с существующими комнатами
		overlap := false
		for _, other := range rooms {
			if newRoom.Intersects(other) {
				overlap = true
				break
			}
		}
		if !overlap {
			rooms = append(rooms, newRoom)
			// Растягиваем стены комнаты в проходимые клетки
			for y := newRoom.Y; y < newRoom.Y+newRoom.Height; y++ {
				for x := newRoom.X; x < newRoom.X+newRoom.Width; x++ {
					gm.Tiles[y][x].Walkable = true
					gm.Tiles[y][x].Symbol = '.'
				}
			}

			// Устанавливаем стартовую позицию игрока в первую комнату
			if len(rooms) == 1 {
				gm.PlayerStartX = newRoom.X + newRoom.Width/2
				gm.PlayerStartY = newRoom.Y + newRoom.Height/2
			}

			// Соединяем текущую комнату с предыдущей коридорами
			if len(rooms) > 1 {
				prevCenterX := previousRoom.X + previousRoom.Width/2
				prevCenterY := previousRoom.Y + previousRoom.Height/2
				currCenterX := newRoom.X + newRoom.Width/2
				currCenterY := newRoom.Y + newRoom.Height/2

				// Случайным образом выбираем порядок соединения
				if rand.Intn(2) == 0 {
					gm.CreateHTunnel(prevCenterX, currCenterX, prevCenterY)
					gm.CreateVTunnel(prevCenterY, currCenterY, currCenterX)
				} else {
					gm.CreateVTunnel(prevCenterY, currCenterY, prevCenterX)
					gm.CreateHTunnel(prevCenterX, currCenterX, currCenterY)
				}
			}

			previousRoom = newRoom
		}
	}

	// Назначаем последнюю комнату как комнату выхода
	if len(rooms) > 0 {
		exitRoom := rooms[len(rooms)-1]
		gm.ExitX = exitRoom.X + rand.Intn(exitRoom.Width)
		gm.ExitY = exitRoom.Y + rand.Intn(exitRoom.Height)
		gm.Tiles[gm.ExitY][gm.ExitX].Symbol = 'E' // Символ выхода
		gm.Tiles[gm.ExitY][gm.ExitX].Walkable = true
	}
}

func (gm *GameMap) CreateHTunnel(x1, x2, y int) {
	for x := min(x1, x2); x <= max(x1, x2); x++ {
		gm.Tiles[y][x].Walkable = true
		gm.Tiles[y][x].Symbol = '.'
	}
}

func (gm *GameMap) CreateVTunnel(y1, y2, x int) {
	for y := min(y1, y2); y <= max(y1, y2); y++ {
		gm.Tiles[y][x].Walkable = true
		gm.Tiles[y][x].Symbol = '.'
	}
}

func (gm *GameMap) Draw(screen *ebiten.Image) {
	for y, row := range gm.Tiles {
		for x, tile := range row {
			var clr color.RGBA
			switch tile.Symbol {
			case '#':
				clr = ColorWhite
			case '.':
				clr = ColorGreen
			case 'E':
				clr = ColorBlue
			}
			text.Draw(screen, string(tile.Symbol), basicfont.Face7x13, x*TileSize, y*TileSize+UIHeight+13, clr)
		}
	}

	// Отрисовка предметов
	for _, item := range gm.Items {
		text.Draw(screen, "!", basicfont.Face7x13, item.X*TileSize, item.Y*TileSize+UIHeight+13, ColorYellow)
	}

	// Отрисовка монстров
	for _, monster := range gm.Monsters {
		monster.Draw(screen)
	}
}

func (gm *GameMap) SpawnMonsters(count int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		for {
			x := rand.Intn(gm.Width)
			y := rand.Intn(gm.Height)
			// Убедимся, что монстр не спавнится на стартовой позиции или на выходе
			if gm.Tiles[y][x].Walkable && !(x == gm.PlayerStartX && y == gm.PlayerStartY) && !(x == gm.ExitX && y == gm.ExitY) && !gm.IsOccupied(x, y) {
				monster := NewMonster(x, y)
				gm.Monsters = append(gm.Monsters, monster)
				break
			}
		}
	}
}

func (gm *GameMap) SpawnItems(count int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		for {
			x := rand.Intn(gm.Width)
			y := rand.Intn(gm.Height)
			// Убедимся, что предмет не спавнится на стартовой позиции или на выходе
			if gm.Tiles[y][x].Walkable && !(x == gm.PlayerStartX && y == gm.PlayerStartY) && !(x == gm.ExitX && y == gm.ExitY) && !gm.IsOccupied(x, y) {
				item := NewItem(x, y)
				gm.Items = append(gm.Items, item)
				break
			}
		}
	}
}

func (gm *GameMap) IsOccupied(x, y int) bool {
	for _, m := range gm.Monsters {
		if m.X == x && m.Y == y {
			return true
		}
	}
	for _, i := range gm.Items {
		if i.X == x && i.Y == y {
			return true
		}
	}
	return false
}

// Вспомогательные функции
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Структура для представления прямоугольника комнаты
type Rectangle struct {
	X, Y, Width, Height int
}

func (r Rectangle) Intersects(other Rectangle) bool {
	return r.X <= other.X+other.Width && r.X+r.Width >= other.X &&
		r.Y <= other.Y+other.Height && r.Y+r.Height >= other.Y
}
