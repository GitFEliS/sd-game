package game

import "testing"

func TestMonsterMove(t *testing.T) {
	gmap := &GameMap{
		Width:  5,
		Height: 5,
		Tiles:  makeTestTiles(5, 5),
	}
	player := NewPlayer(2, 2)
	monster := &Monster{X: 0, Y: 0, Health: 30, Attack: 5}

	// Принудительно задаём поведение монстра, чтобы тест был детерминированным.
	monster.Behavior = &AggressiveBehavior{}

	gmap.Monsters = []*Monster{monster}

	monster.Move(gmap, player)

	// Проверяем, что монстр двинулся к игроку (2,2) из (0,0):
	// Ожидается шаг по диагонали до (1,1).
	if monster.X != 1 || monster.Y != 1 {
		t.Errorf("Expected monster to move towards player, got (%d,%d)", monster.X, monster.Y)
	}
}

func TestMonsterPassive(t *testing.T) {
	gmap := &GameMap{Width: 5, Height: 5, Tiles: makeTestTiles(5, 5)}
	player := NewPlayer(2, 2)
	monster := &Monster{X: 0, Y: 0, Health: 30, Attack: 5, Behavior: &PassiveBehavior{}}
	gmap.Monsters = []*Monster{monster}

	monster.Move(gmap, player)
	// Ожидаем, что монстр не двинется
	if monster.X != 0 || monster.Y != 0 {
		t.Errorf("Expected monster to remain at (0,0), got (%d,%d)", monster.X, monster.Y)
	}
}

func TestMonsterCoward(t *testing.T) {
	gmap := &GameMap{Width: 5, Height: 5, Tiles: makeTestTiles(5, 5)}
	player := NewPlayer(2, 2)
	monster := &Monster{X: 3, Y: 3, Health: 30, Attack: 5, Behavior: &CowardBehavior{}}
	gmap.Monsters = []*Monster{monster}

	monster.Move(gmap, player)
	// Трусливый монстр должен отдаляться, т.е. из (3,3) в сторону от (2,2).
	// Значит должен пойти к (4,4) если та клетка проходима.
	// Проверяем что монстр приблизился к (4,4).
	if monster.X != 4 || monster.Y != 4 {
		t.Errorf("Expected coward monster to move away to (4,4), got (%d,%d)", monster.X, monster.Y)
	}
}
