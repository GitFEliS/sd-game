package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Player struct {
	X, Y            int
	BaseHealth      int
	Health          int
	BaseAttack      int
	Attack          int
	Inventory       []*Item
	EquippedArmor   *Item
	EquippedWeapons []*Item // Максимум 2 оружия
	Level           int
	Exp             int
	NextLevelExp    int
}

func NewPlayer(x, y int) *Player {
	baseHealth := 200
	baseAttack := 20
	return &Player{
		X:               x,
		Y:               y,
		BaseHealth:      baseHealth,
		Health:          baseHealth,
		BaseAttack:      baseAttack,
		Attack:          baseAttack,
		Inventory:       []*Item{},
		EquippedArmor:   nil,
		EquippedWeapons: []*Item{},
		Level:           1,
		Exp:             0,
		NextLevelExp:    100,
	}
}

func (p *Player) Update(g *Game) {
	dx, dy := 0, 0
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		dy = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		dy = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		dx = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		dx = 1
	}

	newX, newY := p.X+dx, p.Y+dy
	if newX >= 0 && newX < g.Map.Width && newY >= 0 && newY < g.Map.Height {
		if g.Map.Tiles[newY][newX].Walkable {
			// Проверка наличия монстра на новой позиции
			for _, monster := range g.Map.Monsters {
				if monster.X == newX && monster.Y == newY {
					p.AttackMonster(monster, g)
					return
				}
			}

			// Проверка наличия предмета на новой позиции
			for i, item := range g.Map.Items {
				if item.X == newX && item.Y == newY {
					p.PickUpItem(g.Map.Items[i], g)
					// Удаление предмета из карты
					g.Map.Items = append(g.Map.Items[:i], g.Map.Items[i+1:]...)
					break
				}
			}

			p.X, p.Y = newX, newY
		}
	}

	// Обработка ввода для открытия инвентаря
	if inpututil.IsKeyJustPressed(ebiten.KeyI) { // Изменили на 'I' для инвентаря
		g.IsInventoryOpen = !g.IsInventoryOpen
		if !g.IsInventoryOpen {
			g.SelectedInventory = 0
		}
	}

	// Проверка уровня
	if p.Exp >= p.NextLevelExp {
		p.LevelUp()
	}
}

func (p *Player) AttackMonster(monster *Monster, g *Game) {
	// Игрок атакует монстра
	monster.Health -= p.Attack
	// Монстр атакует игрока в ответ (симметричный бой только здесь!)
	p.Health -= monster.Attack
	p.Exp += 10 // Опыт за атаку монстра

	// Проверяем, убит ли монстр
	if monster.Health <= 0 {
		// Удаляем монстра из карты
		for i, m := range g.Map.Monsters {
			if m == monster {
				g.Map.Monsters = append(g.Map.Monsters[:i], g.Map.Monsters[i+1:]...)
				break
			}
		}
	}

	// Проверяем, умер ли игрок
	if p.Health <= 0 {
		g.GameOver = true
	}
}

func (p *Player) PickUpItem(item *Item, g *Game) {
	p.Inventory = append(p.Inventory, item)
	// Можно добавить уведомление о подборе предмета
}

func (p *Player) EquipItem(item *Item) {
	slot := item.Type // "Weapon" или "Armor"
	if slot == ItemTypeArmor {
		if p.EquippedArmor != nil {
			// Снимаем текущую броню
			p.UnequipItem(p.EquippedArmor)
		}
		// Экипируем новую броню
		p.EquippedArmor = item
	} else if slot == ItemTypeWeapon {
		if len(p.EquippedWeapons) >= 2 {
			// Автоматически снимаем первое оружие
			p.UnequipItem(p.EquippedWeapons[0])
		}
		// Экипируем новое оружие
		p.EquippedWeapons = append(p.EquippedWeapons, item)
	}
	p.RecalculateStats()
}

func (p *Player) UnequipItem(item *Item) {
	slot := item.Type
	if slot == ItemTypeArmor && p.EquippedArmor == item {
		p.EquippedArmor = nil
	} else if slot == ItemTypeWeapon {
		for i, eqItem := range p.EquippedWeapons {
			if eqItem == item {
				p.EquippedWeapons = append(p.EquippedWeapons[:i], p.EquippedWeapons[i+1:]...)
				break
			}
		}
	}
	p.RecalculateStats()
}

func (p *Player) RecalculateStats() {
	// Сброс базовых характеристик
	p.Health = p.BaseHealth
	p.Attack = p.BaseAttack

	// Применение модификаторов от экипированных предметов
	if p.EquippedArmor != nil {
		for stat, value := range p.EquippedArmor.Modifiers {
			switch stat {
			case "Health":
				p.Health += value
			case "Attack":
				p.Attack += value
			}
		}
	}

	for _, weapon := range p.EquippedWeapons {
		for stat, value := range weapon.Modifiers {
			switch stat {
			case "Health":
				p.Health += value
			case "Attack":
				p.Attack += value
			}
		}
	}
}

func (p *Player) LevelUp() {
	p.Level++
	p.Exp -= p.NextLevelExp
	p.NextLevelExp += 50
	p.BaseHealth += 20
	p.BaseAttack += 5
	p.RecalculateStats()
}

// Метод для проверки, экипирован ли предмет
func (p *Player) IsEquipped(item *Item) bool {
	if p.EquippedArmor == item {
		return true
	}
	for _, weapon := range p.EquippedWeapons {
		if weapon == item {
			return true
		}
	}
	return false
}

func (p *Player) Draw(screen *ebiten.Image) {
	// Отрисовка игрока
	text.Draw(screen, "@", basicfont.Face7x13, p.X*TileSize, p.Y*TileSize+UIHeight+13, ColorWhite)

	// Отображение характеристик игрока
	stats := fmt.Sprintf("Health: %d | Attack: %d | Level: %d | Exp: %d/%d",
		p.Health, p.Attack, p.Level, p.Exp, p.NextLevelExp)
	text.Draw(screen, stats, basicfont.Face7x13, 10, ScreenHeight-10, ColorWhite)

	// Отображение инвентаря
	inventory := "Inventory: "
	for _, item := range p.Inventory {
		inventory += item.Name + " "
	}
	text.Draw(screen, inventory, basicfont.Face7x13, 10, ScreenHeight-25, ColorYellow)

	// Отображение экипированных предметов с суммой модификаторов
	totalArmorHealth := 0
	totalWeaponAttack := 0
	if p.EquippedArmor != nil {
		if val, ok := p.EquippedArmor.Modifiers["Health"]; ok {
			totalArmorHealth += val
		}
	}
	for _, weapon := range p.EquippedWeapons {
		if val, ok := weapon.Modifiers["Attack"]; ok {
			totalWeaponAttack += val
		}
	}

	equipped := "Equipped: "
	if p.EquippedArmor != nil {
		equipped += fmt.Sprintf("Armor(%s) ", p.EquippedArmor.Name)
	} else {
		equipped += "Armor(None) "
	}

	equipped += fmt.Sprintf("Total Armor Health: +%d ", totalArmorHealth)
	equipped += fmt.Sprintf("Weapons Total Attack: +%d", totalWeaponAttack)

	text.Draw(screen, equipped, basicfont.Face7x13, 10, ScreenHeight-40, ColorBlue)
}
