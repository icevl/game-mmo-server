package gameserver

import (
	"fmt"
	"math"
	"server/entity"
	"server/proto/interactpb"
	"server/types"
)

func ActionInteract(world *World, client *types.TCPClient, action *interactpb.Interact) {

	source, err := world.getObject(client.UUID)
	if err != nil {
		fmt.Println("Error getting object from world")
		return
	}

	closes, dist := world.findClosestMapObjectByKind(source, types.ObjectKindTree)
	fmt.Println("Closest object: ", dist)
	if dist <= 1.5 && !closes.IsDead() {
		fmt.Println("Interacting with tree")

		direction := types.Vector3{
			X: closes.Position.X - source.Position.X,
			Y: 0,
			Z: closes.Position.Z - source.Position.Z,
		}
		directionMagnitude := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y + direction.Z*direction.Z)
		normalizedDirection := types.Vector3{
			X: direction.X / directionMagnitude,
			Y: direction.Y / directionMagnitude,
			Z: direction.Z / directionMagnitude,
		}
		angle := math.Atan2(normalizedDirection.X, normalizedDirection.Z) * (180.0 / math.Pi)

		world.transformObjectRotation(source, types.Vector3{X: 0, Y: angle, Z: 0})
		//

		closes.TakeDamage(35)
		world.broadcastSound(closes.Entity.DamageSound, closes.Position, 0.5)

		if !closes.IsDead() {
			world.interactQueue(source)
			return
		}

		world.updateObjectVariation(closes, 1)
		closes.ScheduleRespawn()

	}
}

func ActionInteractWith(world *World, client *types.TCPClient, action *interactpb.InteractWith) {

	source, err := world.getObject(client.UUID)
	if err != nil {
		fmt.Println("Error getting object from world")
		return
	}

	target, err := world.getObject(action.TargetUuid)
	if err != nil {
		fmt.Println("Error getting object from world")
		return
	}

	attackRange := source.GetAttackRange()
	dist := distance(source.Position, target.Position)

	if dist > *attackRange {
		fmt.Println("Target is out of range")
		return

	}

	maxDamage := source.GetAttackMaxDamage()
	if maxDamage == nil || *maxDamage == 0 {
		fmt.Println("Error getting max damage")
		return
	}

	damage := *maxDamage
	isCrit := false

	if damage == *maxDamage {
		isCrit = true
	}

	target.TakeDamage(damage)

	lookAtRotation := source.LookAt(target)
	UpdateTransformRotationChannel <- &types.TransformRotation{Object: source, Rotation: lookAtRotation}

	if target.IsDead() {
		// loot
		world.dropItemOnGround(entity.Pistol, types.Vector3{X: target.Position.X, Y: source.Position.Y, Z: target.Position.Z})
		objectTarget := &target
		world.hideObject(target.UUID)

		DestroyObjectChannel <- &types.DestroyObject{Object: *objectTarget}
		return
	}

	DamageChannel <- &types.Damage{Object: target, Amount: damage, IsCrit: isCrit, HealthCurrent: int32(target.Entity.Health), HealthMax: int32(target.Entity.MaxHealth)}

}
