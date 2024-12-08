package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"math/rand"
	"time"
)

type MonsterBehavior interface {
	Move(gm *GameMap, p *Player, m *Monster)
}

// AggressiveBehavior - монстр всегда двигается к игроку
type AggressiveBehavior struct{}

func (a *AggressiveBehavior) Move(gm *GameMap, p *Player, m *Monster) {
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
	m.tryMove(gm, newX, newY)
}

// PassiveBehavior - монстр просто стоит на месте
type PassiveBehavior struct{}

func (pbh *PassiveBehavior) Move(gm *GameMap, p *Player, m *Monster) {
	// Ничего не делаем
}

// CowardBehavior - монстр старается держаться подальше от игрока
type CowardBehavior struct{}

func (c *CowardBehavior) Move(gm *GameMap, p *Player, m *Monster) {
	dx, dy := 0, 0
	// Логика обратная агрессивной: двигаемся от игрока
	if p.X > m.X {
		dx = -1
	} else if p.X < m.X {
		dx = 1
	}
	if p.Y > m.Y {
		dy = -1
	} else if p.Y < m.Y {
		dy = 1
	}

	newX, newY := m.X+dx, m.Y+dy
	m.tryMove(gm, newX, newY)
}

// ConfusedBehavior - декоратор, оборачивающий другое поведение
// В течение 2 секунд монстр двигается случайно, после чего возвращается к исходному поведению
type ConfusedBehavior struct {
	Inner MonsterBehavior
	Start time.Time
}

func (c *ConfusedBehavior) Move(gm *GameMap, p *Player, m *Monster) {
	// Проверяем, истекло ли 4 секунды с начала конфузии
	if time.Since(c.Start) < 4*time.Second {
		// Конфузия активна - двигаемся случайно
		dx := rand.Intn(3) - 1 // [-1,0,1]
		dy := rand.Intn(3) - 1
		newX, newY := m.X+dx, m.Y+dy
		m.tryMove(gm, newX, newY)
	} else {
		m.Behavior = c.Inner
	}
}

type Monster struct {
	X, Y           int
	Health         int
	Attack         int
	AttackCooldown time.Duration
	LastAttackTime time.Time
	Behavior       MonsterBehavior
}

func NewMonster(x, y int) *Monster {
	health := randInt(20, 60) // Здоровье монстра от 20 до 40
	attack := randInt(5, 10)  // Атака монстра от 5 до 10
	cooldown := time.Second   // Период между атаками

	var behavior MonsterBehavior
	switch rand.Intn(3) {
	case 0:
		behavior = &AggressiveBehavior{}
	case 1:
		behavior = &PassiveBehavior{}
	default:
		behavior = &CowardBehavior{}
	}

	return &Monster{
		X:              x,
		Y:              y,
		Health:         health,
		Attack:         attack,
		AttackCooldown: cooldown,
		LastAttackTime: time.Now(),
		Behavior:       behavior,
	}
}

func (m *Monster) Move(gm *GameMap, p *Player) {
	m.Behavior.Move(gm, p, m)
}

func (m *Monster) tryMove(gm *GameMap, newX, newY int) {
	if newX >= 0 && newX < gm.Width && newY >= 0 && newY < gm.Height {
		if gm.Tiles[newY][newX].Walkable && !gm.IsOccupied(newX, newY) {
			m.X, m.Y = newX, newY
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

func (m *Monster) Confuse() {
	m.Behavior = &ConfusedBehavior{
		Inner: m.Behavior,
		Start: time.Now(),
	}
}
