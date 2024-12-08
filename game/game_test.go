package game

import "testing"

func TestNextLevel(t *testing.T) {
	g := NewGame()
	oldLevel := g.CurrentLevel
	// Имитируем достижение выхода
	g.Player.X = g.Map.ExitX
	g.Player.Y = g.Map.ExitY
	g.Update()

	if g.CurrentLevel == oldLevel {
		t.Errorf("Expected next level, got same level %d", g.CurrentLevel)
	}
}

func TestGameOverOnPlayerDeath(t *testing.T) {
	g := NewGame()
	g.Player.Health = 0
	g.Update()
	if !g.GameOver {
		t.Error("Expected game over when player health <= 0")
	}
}
