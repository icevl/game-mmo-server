package gameserver

import (
	"fmt"
	"server/entity"
	"server/types"

	"github.com/google/uuid"
)

func LoadNPC(object Object) {
	fmt.Println("NPC spawned: ", object.name)

	myUUID, _ := uuid.NewUUID()
	waypoints := [][3]float64{}

	waypoints = append(waypoints, [3]float64{float64(object.position[0]), float64(object.position[2]), float64(object.rotation[1])})
	for _, waypoint := range object.waypoints {
		waypoints = append(waypoints, [3]float64{float64(waypoint[0]), float64(waypoint[2]), 0})
	}

	position := types.Vector3{X: float64(object.position[0]), Y: float64(object.position[1]), Z: float64(object.position[2])}
	rotation := types.Vector3{X: 0, Y: float64(object.rotation[1]), Z: 0}

	npc := &types.GameObject{
		Entity:        entity.EntityFactory(object.name),
		UUID:          myUUID.String(),
		Position:      position,
		Rotation:      rotation,
		PositionSpawn: position,
		RotationSpawn: rotation,
		Waypoints:     waypoints,
		Type:          types.ObjectTypeNPC,
	}

	npc.Entity.Health = npc.Entity.MaxHealth

	if npc.Entity.EquippedItems.RightHand.Type != "" {
		npc.Entity.EquippedItems.RightHand.Clip = npc.Entity.EquippedItems.RightHand.ClipSize
	}

	npc.SetNextTravelTime()

	W.addObject(npc)
	W.updateNeighbors(npc)
}
