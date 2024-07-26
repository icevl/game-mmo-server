package entity

import (
	"math/rand"
	"time"
)

type EquippedItems struct {
	RightHand Entity
	LeftHand  Entity
}

type HumanSlot struct {
	Recipe string
	Color  string
}

type HumanCharacter struct {
	Gender string
	Slots  map[string]HumanSlot
}

type EntityType = string

const (
	TypePistol EntityType = "pistol"
	TypeAxe    EntityType = "axe"
)

type Entity struct {
	Name            string
	RespawnInterval int // in seconds
	InternalName    string
	Resource        string
	Type            EntityType
	MaxHealth       int32
	Health          int32
	InteractChance  float32 // 0 - 100
	DamageSound     string
	CanAgro         bool

	EquippedItems *EquippedItems

	Variation    string // basic, dragon, etc
	ClipSize     int    // for pistols, guns, etc
	Clip         int    // current clip size
	ReloadTime   float32
	AttackDamage int32
	AttackRange  float32
	AttackSpeed  float32
	AttackRadius float32

	ReloadFinishTime *time.Time // Time to finish reload

	Speed          float32
	HumanCharacter *HumanCharacter
}

func EntityFactory(internalName string) Entity {
	switch internalName {
	case "tree":
		return Tree
	case "adam":
		return Bandit
	case "cyber_woman":
		return CyberWoman
	}

	return Entity{}
}

func (e *Entity) CanInteract() bool {
	if e.InteractChance == 0 {
		return true
	}
	chance := rand.Float32() * 100
	return chance <= e.InteractChance
}

func (e *Entity) GetInteractAnimation() string {
	var rightHand Entity

	if e.EquippedItems.RightHand.Type != "" {
		rightHand = e.EquippedItems.RightHand
	}

	switch rightHand.Type {
	case TypePistol:
		return "PistolShoot"
	case TypeAxe:
		return "AttackAxe"
	default:
		return ""
	}
}

func (e *Entity) StartReloadWeapon() {
	if e.EquippedItems.RightHand.ReloadTime == 0 {
		return
	}

	e.EquippedItems.RightHand.Clip = e.EquippedItems.RightHand.ClipSize
	reloadFinishedAt := time.Now().Add(time.Duration(e.EquippedItems.RightHand.ReloadTime*1000) * time.Millisecond)
	e.EquippedItems.RightHand.ReloadFinishTime = &reloadFinishedAt
}

func (e *Entity) IsReloadWeaponInProgress() bool {
	if e.EquippedItems.RightHand.ReloadFinishTime == nil {
		return false
	}

	return time.Now().Before(*e.EquippedItems.RightHand.ReloadFinishTime)
}

var axes = []Entity{BasicAxe, DragonAxe}
var pistols = []Entity{Pistol}
var itemsMap map[string]Entity = make(map[string]Entity)

func init() {
	var items []Entity

	items = append(items, axes...)
	items = append(items, pistols...)

	for _, item := range items {
		itemsMap[item.InternalName] = item
	}
}

func GetItem(name string) Entity {
	if item, ok := itemsMap[name]; ok {
		return item
	}

	return Entity{}
}
