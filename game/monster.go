package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"time"
)

type Monster struct {
	X, Y           int
	Health         int
	Attack         int
	AttackCooldown time.Duration
	LastAttackTime time.Time
}

func NewMonster(x, y int) *Monster {
	health := randInt(20, 60) // Здоровье монстра от 20 до 40
	attack := randInt(5, 70)  // Атака монстра от 5 до 10
	cooldown := time.Second   // Период между атаками
	return &Monster{
		X:              x,
		Y:              y,
		Health:         health,
		Attack:         attack,
		AttackCooldown: cooldown,
		LastAttackTime: time.Now(),
	}
}

func (m *Monster) Move(gm *GameMap, p *Player) {
	// Простая логика: перемещение в сторону игрока
	dx, dy := 0, 0
	if p.X > m.X {
		dx = 1
	} else if p.X < m.X {
		dx = -1
	}
	if p.Y > m.Y {
		dy = 1
	} else if p.Y < m.Y {
		dy = -1
	}

	newX, newY := m.X+dx, m.Y+dy
	if newX >= 0 && newX < gm.Width && newY >= 0 && newY < gm.Height {
		if gm.Tiles[newY][newX].Walkable && !gm.IsOccupied(newX, newY) {
			m.X, m.Y = newX, newY
		}
	}

	// Проверка столкновения с игроком и атака с учетом таймера
	if m.X == p.X && m.Y == p.Y {
		now := time.Now()
		if now.Sub(m.LastAttackTime) >= m.AttackCooldown {
			p.Health -= m.Attack
			m.LastAttackTime = now
			fmt.Printf("Player HP: %d\n", p.Health)
		}
	}
}

func (m *Monster) Draw(screen *ebiten.Image) {
	// Отрисовка монстра
	text.Draw(screen, "M", basicfont.Face7x13, m.X*TileSize, m.Y*TileSize+UIHeight+13, ColorRed)

	// Отображение здоровья монстра
	healthMsg := fmt.Sprintf("HP: %d, DMG: %d", m.Health, m.Attack)
	text.Draw(screen, healthMsg, basicfont.Face7x13, m.X*TileSize, m.Y*TileSize+UIHeight-5, ColorYellow)
}
