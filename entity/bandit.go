package entity

var Bandit = Entity{
	Name:            "Bandit",
	InternalName:    "bandit",
	MaxHealth:       100,
	RespawnInterval: 60,
	CanAgro:         true,
	Speed:           2,
	HumanCharacter: &HumanCharacter{
		Gender: "male",
		Slots: map[string]HumanSlot{
			"Hair":  {Recipe: "MilCut", Color: "#FFFFFF"},
			"Beard": {Recipe: "MaleBeard1", Color: "#FFFFFF"},
			"Legs":  {Recipe: "MaleSweatPants_Recipe", Color: "#FF0000"},
			"Chest": {Recipe: "MaleShirt2", Color: "#FF0000"},
		},
	},
	// EquippedItems: &EquippedItems{
	// 	RightHand: BasicAxe,
	// },

	EquippedItems: &EquippedItems{
		RightHand: BasicAxe,
	},
}
