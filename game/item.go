package game

type Item struct {
	Name      string
	Type      string
	Modifiers map[string]int
	X, Y      int
}

func (i *Item) Use(p *Player) {
	for stat, value := range i.Modifiers {
		switch stat {
		case "Health":
			p.Health += value
		case "Attack":
			p.Attack += value
		}
	}
}
