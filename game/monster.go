package game

import (
	"math/rand"
)

type Monster struct {
	X, Y    int
	Health  int
	AttackP int
}

func (m *Monster) Move(gm *GameMap) {
	dx := rand.Intn(3) - 1
	dy := rand.Intn(3) - 1
	newX, newY := m.X+dx, m.Y+dy
	if newX >= 0 && newX < gm.Width && newY >= 0 && newY < gm.Height {
		if gm.Tiles[newY][newX].Walkable {
			m.X, m.Y = newX, newY
		}
	}
}

func (m *Monster) Attack(p *Player) {
	p.Health -= m.AttackP
}
