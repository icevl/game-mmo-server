package events

import (
	pbglobal "server/proto"
	"server/proto/actionpb"
	"server/proto/animationpb"
	"server/proto/interactpb"
	"server/proto/objectpb"
	"server/types"
)

func GetDestroyObjectEventPayload(uuid string) *actionpb.Action {
	return &actionpb.Action{
		Action: &actionpb.Action_DestroyObject{
			DestroyObject: &objectpb.DestroyObject{
				UUID: uuid,
			},
		},
	}
}

func GetAnimationEventPayload(uuid string, name string, speed float32, isStop bool) *actionpb.Action {
	return &actionpb.Action{
		Action: &actionpb.Action_Animation{
			Animation: &animationpb.Animation{
				UUID:   uuid,
				Name:   name,
				Speed:  speed,
				IsStop: isStop,
			},
		},
	}
}

func GetDamagePayload(uuid string, amount int32, isCrit bool, currentHealth, maxHealth int32) *actionpb.Action {
	return &actionpb.Action{
		Action: &actionpb.Action_Damage{
			Damage: &objectpb.Damage{
				UUID:          uuid,
				Amount:        amount,
				IsCrit:        isCrit,
				HealthCurrent: currentHealth,
				HealthMax:     maxHealth,
			},
		},
	}
}

func GetObjectEventPayload(object *types.GameObject, options *types.EventPayloadOptions) *actionpb.Action {
	return &actionpb.Action{
		Action: &actionpb.Action_Object{
			Object: GetObjectEvent(object, options),
		},
	}
}

func GetInteractQueuePayload(object *types.GameObject) *actionpb.Action {
	return &actionpb.Action{
		Action: &actionpb.Action_InteractQueue{
			InteractQueue: &interactpb.InteractQueue{
				UUID: object.UUID,
			},
		},
	}
}

func GetObjectEvent(object *types.GameObject, options *types.EventPayloadOptions) *objectpb.Object {

	msg := &objectpb.Object{
		UUID:      object.UUID,
		Name:      object.Name,
		Resource:  object.Resource,
		Speed:     object.Speed,
		IsSelf:    *&options.IsSelf,
		Type:      object.Type,
		Variation: object.Entity.Variation,
		Position: &pbglobal.Vector3M{
			X: float32(object.Position.X),
			Y: float32(object.Position.Y),
			Z: float32(object.Position.Z),
		},
		Rotation: &pbglobal.Vector3M{
			X: float32(object.Rotation.X),
			Y: float32(object.Rotation.Y),
			Z: float32(object.Rotation.Z),
		},
	}

	if object.EquippedItems != nil {
		msg.EquippedItems = &objectpb.EquippedItems{}

		if object.EquippedItems.RightHand.Resource != "" {
			msg.EquippedItems.RightHand = &objectpb.ObjectWithVariation{
				Resource:     object.EquippedItems.RightHand.Resource,
				Variation:    object.EquippedItems.RightHand.Variation,
				AttackRadius: object.EquippedItems.RightHand.AttackRadius,
			}
		}

		if object.EquippedItems.LeftHand.Resource != "" {
			msg.EquippedItems.LeftHand = &objectpb.ObjectWithVariation{
				Resource:     object.EquippedItems.LeftHand.Resource,
				Variation:    object.EquippedItems.LeftHand.Variation,
				AttackRadius: object.EquippedItems.LeftHand.AttackRadius,
			}
		}

	}

	if object.HumanCharacter != nil {
		slots := make(map[string]*objectpb.HumanSlot)

		for key, slot := range object.HumanCharacter.Slots {
			slots[key] = &objectpb.HumanSlot{
				Recipe: slot.Recipe,
				Color:  slot.Color,
			}
		}

		msg.HumanCharacter = &objectpb.HumanCharacter{
			Gender: object.HumanCharacter.Gender,
			Slots:  slots,
		}
	}

	return msg
}

func GetNetworkStatePayload(object *types.GameObject) *actionpb.Action {
	return &actionpb.Action{
		Action: &actionpb.Action_ObjectState{
			ObjectState: &objectpb.ObjectState{
				UUID:         object.UUID,
				VariantIndex: object.VariationIndex,
			},
		},
	}
}
