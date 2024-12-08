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
	gmap.Monsters = []*Monster{monster}

	monster.Move(gmap, player)

	// Монстр должен двинуться в сторону игрока (2,2).
	// Из (0,0) следующий шаг к игроку будет (1,1) (по диагонали), если логика так работает.
	// Если монстр двигается на одну клетку по x и одну по y, тогда проверим новые координаты.
	// По коду: dx = 1, dy = 1 => newX=1, newY=1
	if monster.X != 1 || monster.Y != 1 {
		t.Errorf("Expected monster to move towards player, got (%d,%d)", monster.X, monster.Y)
	}
}
