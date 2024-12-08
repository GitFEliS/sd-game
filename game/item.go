package game

import (
	"fmt"
	rand2 "math/rand/v2"
)

type Item struct {
	Name      string
	Type      string // "Weapon" или "Armor"
	Modifiers map[string]int
	X, Y      int
}

func NewItem(x, y int) *Item {
	itemType := ItemTypeWeapon
	if rand2.Int32()%2 == 0 {
		itemType = ItemTypeArmor
	}

	modifiers := make(map[string]int)
	if itemType == ItemTypeWeapon {
		modifiers["Attack"] = randInt(5, 10) // Увеличение атаки на 5-10
	} else if itemType == ItemTypeArmor {
		modifiers["Health"] = randInt(10, 30) // Увеличение здоровья на 10-30
	}

	name := itemType
	if itemType == ItemTypeWeapon {
		name += fmt.Sprintf(" +%d", modifiers["Attack"])
	} else if itemType == ItemTypeArmor {
		name += fmt.Sprintf(" +%d", modifiers["Health"])
	}

	return &Item{
		Name:      name,
		Type:      itemType,
		Modifiers: modifiers,
		X:         x,
		Y:         y,
	}
}

// Вспомогательная функция для генерации случайных чисел
func randInt(min, max int) int {
	return min + int(rand2.Int32()%int32(max-min+1))
}
