package entity

var Pistol = Entity{
	Name:         "Colt",
	InternalName: "pistol",
	Type:         TypePistol,
	ClipSize:     100, //7,
	Health:       1000,
	MaxHealth:    1000,
	ReloadTime:   2.5,
	AttackDamage: 1,
	AttackRange:  5,
	AttackSpeed:  1.15,
	AttackRadius: 0.6,
	Resource:     "Weapon/Gun/Pistol",
}
