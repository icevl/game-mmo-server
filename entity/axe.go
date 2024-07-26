package entity

var BasicAxe = Entity{
	Name:         "Basic Axe",
	InternalName: "basic_axe",
	Type:         TypeAxe,
	Health:       100,
	MaxHealth:    100,
	AttackDamage: 1,
	AttackRange:  1.5,
	AttackSpeed:  2,
	AttackRadius: 0.6,
	Resource:     "Weapon/BasicAxe",
}

var DragonAxe = Entity{
	Name:         "Dragon Axe",
	InternalName: "dragon_axe",
	Type:         TypeAxe,
	Health:       1000,
	MaxHealth:    1000,
	AttackDamage: 50,
	AttackRange:  1.5,
	AttackSpeed:  2,
	AttackRadius: 0.6,
	Resource:     "Weapon/BasicAxe",
	Variation:    "dragon",
}
