package entity

var CyberWoman = Entity{
	Name:            "Cyber Woman",
	InternalName:    "cyber_woman",
	MaxHealth:       100,
	RespawnInterval: 60,
	CanAgro:         true,

	Speed:    5,
	Resource: "Characters/WomanCyber",

	EquippedItems: &EquippedItems{
		RightHand: BasicAxe,
	},
}
