package game

import "testing"

func TestIsOccupied(t *testing.T) {
	gm := &GameMap{
		Width:  5,
		Height: 5,
		Tiles:  makeTestTiles(5, 5),
		Monsters: []*Monster{
			{X: 2, Y: 2},
		},
		Items: []*Item{
			{X: 1, Y: 1, Name: "TestItem"},
		},
	}

	if !gm.IsOccupied(2, 2) {
		t.Error("Expected occupied by monster at (2,2)")
	}
	if !gm.IsOccupied(1, 1) {
		t.Error("Expected occupied by item at (1,1)")
	}
	if gm.IsOccupied(0, 0) {
		t.Error("Expected not occupied at (0,0)")
	}
}
