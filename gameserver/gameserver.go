package gameserver

import (
	"fmt"
	"server/entity"
	"server/types"
	"time"
)

var W *World

func StartGameServer() {
	fmt.Println("Starting game server")
	ch := make(chan int)

	level, _ := LoadLevel()

	W = NewWorld(float64(level.TerrainData.size[0]))

	for _, teleport := range level.Teleports {
		W.teleports = append(W.teleports, &teleport)
	}

	for _, object := range level.Objects {

		if object.isNPC() {
			LoadNPC(object)
			continue
		}

		uuid := fmt.Sprintf("object-%d", object.uid)

		levelObject := &types.GameObject{
			Entity:         entity.EntityFactory(object.name),
			Kind:           object.kind,
			UUID:           uuid,
			VariationIndex: object.variationIndex,
			Position:       types.Vector3{X: float64(object.position[0]), Y: float64(object.position[1]), Z: float64(object.position[2])},
			Rotation:       types.Vector3{X: float64(object.rotation[0]), Y: float64(object.rotation[1]), Z: float64(object.rotation[2])},
			Type:           types.ObjectTypeVariantMapObject,
		}

		W.addObject(levelObject)
	}

	go StartUDPServer()
	go StartTCPServer()

	go globalTicker()

	go ProcessMovementUpdates()
	go ProcessTeleportObjectUpdates()
	go ProcessGameObjectVariationsUpdates()
	go ProcessSoundBroadcast()
	go ProcessTransformRotationUpdates()
	go ProcessAnimationUpdates()
	go ProcessDamage()
	go ProcessObjectDestroy()
	go ProcessSpawnObject()
	go ProcessInteractQueue()

	ch <- 1

}

func globalTicker() {

	// Run global ticker 25 times per second
	ticker := time.NewTicker(40 * time.Millisecond)
	for range ticker.C {
		W.npcRespawnTick()
		W.npcWalkTick()
		W.npcAttackTick()
		W.mapObjectVariationTick()
		W.mapObjectDestroyTick()
		// currentTime := time.Now()
		// fmt.Println("Tick: ", currentTime)
	}
}
