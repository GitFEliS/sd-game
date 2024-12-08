package game

import (
	"testing"
)

// Тест на поднятие предмета
func TestPlayerPickUpItem(t *testing.T) {
	// Подготовка тестовых данных
	gmap := &GameMap{
		Width:  10,
		Height: 10,
		Tiles:  makeTestTiles(10, 10),
		Items:  []*Item{{X: 1, Y: 1, Name: "TestItem", Type: ItemTypeWeapon, Modifiers: map[string]int{"Attack": 5}}},
	}
	player := NewPlayer(0, 0)
	g := &Game{Map: gmap, Player: player}

	// Помещаем игрока рядом с предметом
	player.X, player.Y = 0, 0
	// Перемещаем игрока на клетку с предметом
	movePlayer(g, 1, 1)
	// Проверяем, поднялся ли предмет
	if len(player.Inventory) != 1 {
		t.Errorf("Expected 1 item in inventory, got %d", len(player.Inventory))
	}
	if len(gmap.Items) != 0 {
		t.Errorf("Expected no items on map after pickup, got %d", len(gmap.Items))
	}
}

// Тест на экипировку предмета
func TestPlayerEquipItem(t *testing.T) {
	player := NewPlayer(0, 0)
	weapon := &Item{Name: "Sword", Type: ItemTypeWeapon, Modifiers: map[string]int{"Attack": 5}}
	player.Inventory = append(player.Inventory, weapon)

	player.EquipItem(weapon)
	// Проверяем, применились ли модификаторы
	expectedAttack := player.BaseAttack + 5
	if player.Attack != expectedAttack {
		t.Errorf("Expected attack %d, got %d", expectedAttack, player.Attack)
	}
	if !player.IsEquipped(weapon) {
		t.Error("Weapon should be equipped")
	}
}

// Тест на атаку монстра
func TestPlayerAttackMonster(t *testing.T) {
	player := NewPlayer(0, 0)
	player.Health = 100
	player.Attack = 10
	monster := &Monster{
		X: 0, Y: 0,
		Health: 50,
		Attack: 5,
	}

	g := &Game{
		Player: player,
		Map:    &GameMap{Monsters: []*Monster{monster}},
	}

	// Игрок атакует монстра
	player.AttackMonster(monster, g)

	// После атаки у монстра должно стать меньше здоровья
	if monster.Health != 40 { // 50 - 10 = 40
		t.Errorf("Expected monster health 40, got %d", monster.Health)
	}
	// Игрок получает ответный удар
	if player.Health != 95 { // 100 - 5 = 95
		t.Errorf("Expected player health 95, got %d", player.Health)
	}
	// Игрок получает опыт
	if player.Exp != 10 {
		t.Errorf("Expected player exp 10, got %d", player.Exp)
	}
}

// Вспомогательные функции для тестов

// makeTestTiles создает двумерный массив клеток Tile, все walkable = true
func makeTestTiles(w, h int) [][]Tile {
	tiles := make([][]Tile, h)
	for y := 0; y < h; y++ {
		tiles[y] = make([]Tile, w)
		for x := 0; x < w; x++ {
			tiles[y][x] = Tile{
				Walkable: true,
				Symbol:   '.',
			}
		}
	}
	return tiles
}

// movePlayer - имитация хода игрока для теста.
// Просто меняет координаты игрока, не зависит от реальной логики ввода.
func movePlayer(g *Game, newX, newY int) {
	// Проверяем можно ли идти
	if newX >= 0 && newX < g.Map.Width && newY >= 0 && newY < g.Map.Height && g.Map.Tiles[newY][newX].Walkable {
		// Проверка на предмет
		for i, item := range g.Map.Items {
			if item.X == newX && item.Y == newY {
				g.Player.PickUpItem(item, g)
				g.Map.Items = append(g.Map.Items[:i], g.Map.Items[i+1:]...)
				break
			}
		}
		// Проверка на монстра
		for _, monster := range g.Map.Monsters {
			if monster.X == newX && monster.Y == newY {
				g.Player.AttackMonster(monster, g)
				return
			}
		}
		g.Player.X, g.Player.Y = newX, newY
	}
}
