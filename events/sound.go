package events

import (
	"server/proto"
	"server/proto/actionpb"
	"server/proto/soundpb"
	"server/types"
)

func GetPlaySoundEventPayload(resource string, position types.Vector3, volume float32) *actionpb.Action {
	msg := &soundpb.PlaySound{
		Resource: resource,
		Position: &proto.Vector3M{X: float32(position.X), Y: float32(position.Y), Z: float32(position.Z)},
		Volume:   volume,
	}

	return &actionpb.Action{
		Action: &actionpb.Action_PlaySound{
			PlaySound: msg,
		},
	}
}
