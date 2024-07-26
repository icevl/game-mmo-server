package events

import (
	pbglobal "server/proto"
	"server/proto/actionpb"
	"server/proto/transformpb"
	"server/types"
)

func GetTransformRotationEventPayload(uuid string, rotation types.Vector3) *actionpb.Action {
	msg := &transformpb.TransformRotation{
		UUID: uuid,
		Rotation: &pbglobal.Vector3M{
			X: float32(rotation.X),
			Y: float32(rotation.Y),
			Z: float32(rotation.Z),
		},
	}

	return &actionpb.Action{
		Action: &actionpb.Action_TransformRotation{
			TransformRotation: msg,
		},
	}
}

func GetTeleportEventPayload(uuid string, position, rotation types.Vector3) *actionpb.Action {
	return &actionpb.Action{
		Action: &actionpb.Action_Teleport{
			Teleport: &transformpb.Teleport{
				UUID: uuid,
				Position: &pbglobal.Vector3M{
					X: float32(position.X),
					Y: float32(position.Y),
					Z: float32(position.Z),
				},
				Rotation: &pbglobal.Vector3M{
					X: float32(rotation.X),
					Y: float32(rotation.Y),
					Z: float32(rotation.Z),
				},
			},
		},
	}
}
