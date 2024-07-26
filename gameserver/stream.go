package gameserver

import (
	"fmt"
	"server/events"
	"server/proto"
	"server/proto/actionpb"
	"server/proto/transformpb"
	"server/types"
)

const bufferSize = 1024

var UpdateMovementChannel = make(chan *types.GameObject, bufferSize)
var UpdateTransformRotationChannel = make(chan *types.TransformRotation, bufferSize)
var UpdateGameObjectVariationChannel = make(chan *types.GameObjectVariation, bufferSize)
var BroadcastSoundChannel = make(chan *types.BroadcastSound, bufferSize)
var UpdateAnimationChannel = make(chan *types.Animation, bufferSize)
var DamageChannel = make(chan *types.Damage, bufferSize)
var DestroyObjectChannel = make(chan *types.DestroyObject, bufferSize)
var SpawnObjectChannel = make(chan *types.SpawnObject, bufferSize)
var TeleportObjectChannel = make(chan *types.TeleportObject, bufferSize)
var InteractQueue = make(chan *types.InteractQueue, bufferSize)

func ProcessMovementUpdates() {
	for obj := range UpdateMovementChannel {
		W.onWalkUpdates(obj)

		msg := &transformpb.Transform{
			UUID:     obj.UUID,
			Speed:    obj.Speed,
			Position: &proto.Vector3M{X: float32(obj.Position.X), Y: float32(obj.Position.Y), Z: float32(obj.Position.Z)},
			Rotation: &proto.Vector3M{X: float32(obj.Rotation.X), Y: float32(obj.Rotation.Y), Z: float32(obj.Rotation.Z)},
		}

		for _, player := range obj.GetPlayersNearby() {
			if player.Type != types.ObjectTypePlayer {
				continue
			}

			UDPState.sendToClient(player.UUID, &actionpb.Action{
				Action: &actionpb.Action_Transform{
					Transform: msg,
				},
			})
		}
	}
}

func ProcessGameObjectVariationsUpdates() {
	for request := range UpdateGameObjectVariationChannel {
		msg := events.GetNetworkStatePayload(request.Object)

		for _, player := range request.Object.GetPlayersNearby() {
			TCPState.sendToClient(player.UUID, msg)
		}
	}
}

func ProcessDamage() {
	for request := range DamageChannel {
		msg := events.GetDamagePayload(request.Object.UUID, request.Amount, request.IsCrit, request.HealthCurrent, request.HealthMax)
		if request.Object.Type == types.ObjectTypePlayer {
			TCPState.sendToClient(request.Object.UUID, msg)
		}

		for _, player := range request.Object.GetPlayersNearby() {
			TCPState.sendToClient(player.UUID, msg)
		}
	}
}

func ProcessObjectDestroy() {
	for request := range DestroyObjectChannel {
		msg := events.GetDestroyObjectEventPayload(request.Object.UUID)

		for _, player := range request.Object.GetPlayersNearby() {
			TCPState.sendToClient(player.UUID, msg)
		}
	}
}

func ProcessSpawnObject() {
	for request := range SpawnObjectChannel {
		msg := events.GetObjectEventPayload(request.Object, &types.EventPayloadOptions{IsSelf: false})

		for _, player := range request.Object.GetPlayersNearby() {
			TCPState.sendToClient(player.UUID, msg)
		}
	}
}

func ProcessInteractQueue() {
	for request := range InteractQueue {
		msg := events.GetInteractQueuePayload(request.Object)
		TCPState.sendToClient(request.Object.UUID, msg)
	}
}

func ProcessAnimationUpdates() {
	for request := range UpdateAnimationChannel {
		msg := events.GetAnimationEventPayload(request.Object.UUID, request.Name, request.Speed, request.IsStop)
		fmt.Println("Sending animation update", request.Object.GetPlayersNearby())
		for _, player := range request.Object.GetPlayersNearby() {
			TCPState.sendToClient(player.UUID, msg)
		}
	}
}

func ProcessTransformRotationUpdates() {
	for request := range UpdateTransformRotationChannel {

		players := make([]*types.GameObject, 0)
		players = append(players, request.Object.GetPlayersNearby()...)
		players = append(players, request.Object)

		for _, player := range players {
			msg := events.GetTransformRotationEventPayload(request.Object.UUID, request.Rotation)
			UDPState.sendToClient(player.UUID, msg)
		}
	}
}

func ProcessTeleportObjectUpdates() {
	for request := range TeleportObjectChannel {

		players := make([]*types.GameObject, 0)
		players = append(players, request.Object.GetPlayersNearby()...)
		players = append(players, request.Object)

		for _, player := range players {
			msg := events.GetTeleportEventPayload(request.Object.UUID, request.Position, request.Rotation)
			TCPState.sendToClient(player.UUID, msg)
		}
	}
}

func ProcessSoundBroadcast() {
	for request := range BroadcastSoundChannel {
		listeners := W.getPlayersByPosition(request.Position, 50)
		msg := events.GetPlaySoundEventPayload(request.Resource, request.Position, request.Volume)

		for _, listener := range listeners {
			TCPState.sendToClient(listener.UUID, msg)
		}
	}
}

func findChanges(oldSlice, newSlice []*types.GameObject) (added, removed []*types.GameObject) {
	oldMap := make(map[string]*types.GameObject, len(oldSlice))
	for _, p := range oldSlice {
		oldMap[p.UUID] = p
	}

	newMap := make(map[string]*types.GameObject, len(newSlice))
	for _, p := range newSlice {
		newMap[p.UUID] = p
	}

	for uuid, p := range oldMap {
		if _, ok := newMap[uuid]; !ok {
			removed = append(removed, p)
		}
	}

	for uuid, p := range newMap {
		if _, ok := oldMap[uuid]; !ok {
			added = append(added, p)
		}
	}

	return added, removed
}
