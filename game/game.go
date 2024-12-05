package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
)

const (
	levels = 5
)

type Game struct {
	Player            *Player
	Map               *GameMap
	CurrentLevel      int
	GameOver          bool
	IsInventoryOpen   bool
	SelectedInventory int
}

func NewGame() *Game {
	currentLevel := 1
	gm := NewGameMap(currentLevel)
	player := NewPlayer(gm.PlayerStartX, gm.PlayerStartY)
	return &Game{
		Player:            player,
		Map:               gm,
		CurrentLevel:      currentLevel,
		GameOver:          false,
		IsInventoryOpen:   false,
		SelectedInventory: 0,
	}
}

func (g *Game) Update() error {
	if g.GameOver {
		return nil
	}

	if g.Player.Health <= 0 {
		g.GameOver = true
		return nil
	}

	if g.IsInventoryOpen {
		g.HandleInventoryInput()
		return nil
	}

	g.Player.Update(g)

	// Обновление монстров
	for _, monster := range g.Map.Monsters {
		monster.Move(g.Map, g.Player)
	}

	// Проверка достижения выхода
	if g.Player.X == g.Map.ExitX && g.Player.Y == g.Map.ExitY {
		g.NextLevel()
	}

	return nil
}

func (g *Game) HandleInventoryInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // Закрытие инвентаря
		g.IsInventoryOpen = false
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) { // Навигация вверх
		if g.SelectedInventory > 0 {
			g.SelectedInventory--
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) { // Навигация вниз
		if g.SelectedInventory < len(g.Player.Inventory)-1 {
			g.SelectedInventory++
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) { // Экипировка предмета
		if len(g.Player.Inventory) > 0 {
			selectedItem := g.Player.Inventory[g.SelectedInventory]
			if selectedItem.Type == ItemTypeWeapon || selectedItem.Type == ItemTypeArmor {
				// Проверяем, не экипирован ли уже предмет
				if !g.Player.IsEquipped(selectedItem) {
					g.Player.EquipItem(selectedItem)
				}
			}
		}
	}

	// Возможность снять экипированный предмет с помощью клавиши 'U'
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		if len(g.Player.Inventory) > 0 {
			selectedItem := g.Player.Inventory[g.SelectedInventory]
			// Проверяем, экипирован ли предмет
			if g.Player.IsEquipped(selectedItem) {
				g.Player.UnequipItem(selectedItem)
			}
		}
	}
}

func (g *Game) NextLevel() {
	g.CurrentLevel++
	if g.CurrentLevel > levels {
		g.GameOver = true
		return
	}
	g.Map = NewGameMap(g.CurrentLevel)
	g.Player.X = g.Map.PlayerStartX
	g.Player.Y = g.Map.PlayerStartY
	g.Player.LevelUp()
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.GameOver {
		var msg string
		if g.Player.Health <= 0 {
			msg = "Game Over! You Died."
		} else {
			msg = "Congratulations! You Won."
		}
		text.Draw(screen, msg, basicfont.Face7x13, ScreenWidth/2-100, ScreenHeight/2, ColorRed)
		return
	}

	g.Map.Draw(screen)
	g.Player.Draw(screen)

	// Отображение текущего уровня
	levelMsg := fmt.Sprintf("Level: %d", g.CurrentLevel)
	text.Draw(screen, levelMsg, basicfont.Face7x13, ScreenWidth-100, ScreenHeight-10, ColorYellow)

	// Отображение подсказок
	g.DrawHints(screen)

	// Отображение состояния инвентаря
	if g.IsInventoryOpen {
		g.DrawInventory(screen)
	}
}

func (g *Game) DrawHints(screen *ebiten.Image) {
	// Подсказки отображаются в верхней области UI
	hints := []string{
		"Controls:",
		"W/A/S/D - Move",
		"I - Open/Close Inventory",
		"Arrow Keys - Navigate Inventory",
		"Enter - Equip Item",
		"U - Unequip Item",
		"Esc - Close Inventory",
		"Equip Limit: 1 Armor, 2 Weapons",
		"EXP per Monster: 10",
		"Reach the blue 'E' to complete the level.",
	}

	y := 40 // Начальная позиция внутри UI области (UIHeight=100)
	for _, hint := range hints {
		text.Draw(screen, hint, basicfont.Face7x13, 10, y, ColorGreen)
		y += 15
	}
}

func (g *Game) DrawInventory(screen *ebiten.Image) {
	// Рисуем полупрозрачный фон для инвентаря
	for y := 150; y < 650; y++ { // Изменена граница до 650 для соответствия ScreenHeight=800
		for x := 150; x < 650; x++ {
			screen.Set(x, y, color.RGBA{0, 0, 0, 200}) // Полупрозрачный черный фон
		}
	}

	// Заголовок
	text.Draw(screen, "Inventory:", basicfont.Face7x13, 160, 180, ColorWhite)

	// Список предметов
	for i, item := range g.Player.Inventory {
		colorToUse := ColorWhite
		if i == g.SelectedInventory {
			colorToUse = ColorYellow
		}

		equippedStatus := "Unequipped"
		if g.Player.IsEquipped(item) {
			equippedStatus = "Equipped"
		}

		text.Draw(screen, fmt.Sprintf("%d. %s [%s]", i+1, item.Name, equippedStatus), basicfont.Face7x13, 160, 200+20*i, colorToUse)
	}

	// Инструкция
	text.Draw(screen, "Use Arrow Keys to navigate, Enter to equip, U to unequip, Esc to close.", basicfont.Face7x13, 160, 600, ColorGreen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
