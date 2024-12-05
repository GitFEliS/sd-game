package game

import "testing"

func TestNewItem(t *testing.T) {
	item := NewItem(3, 3)
	if item.X != 3 || item.Y != 3 {
		t.Errorf("Item coordinates not set properly, got (%d,%d)", item.X, item.Y)
	}
	if item.Type != ItemTypeWeapon && item.Type != ItemTypeArmor {
		t.Error("Item type should be either Weapon or Armor")
	}
	if len(item.Modifiers) == 0 {
		t.Error("Expected item to have modifiers")
	}
}
