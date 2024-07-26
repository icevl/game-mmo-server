package gameserver

import (
	"errors"
	"fmt"
	"math"
	"server/entity"
	"server/types"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

const AREA_OF_INTEREST float64 = 30 // 400

type World struct {
	sync.RWMutex
	Octree    *types.Octree
	objects   map[string]*types.GameObject
	teleports []*LevelTeleport
}

type LookedAtObject struct {
	GameObject *types.GameObject
	Distance   float64
}

func NewWorld(size float64) *World {
	oct := types.CreateOctree(
		types.Vector3f{-size, -size, -size},
		types.Vector3f{size, size, size},
	)
	return &World{
		Octree:  oct,
		objects: make(map[string]*types.GameObject),
	}
}

func (w *World) addObject(obj *types.GameObject) {
	w.Lock()
	defer w.Unlock()

	node := w.Octree.Add(obj, types.Vector3f{obj.Position.X, obj.Position.Y, obj.Position.Z})
	obj.Node = node
	w.objects[obj.UUID] = obj
}

func (w *World) getObject(uuid string) (*types.GameObject, error) {
	w.Lock()
	defer w.Unlock()

	var object = w.objects[uuid]

	if object == nil {
		return nil, errors.New("object not found")
	}

	return object, nil
}

func (w *World) removeObject(uuid string) {
	w.Lock()

	obj, ok := w.objects[uuid]
	if !ok {
		w.Unlock()
		return
	}

	neighbors := obj.Neighbors
	delete(w.objects, obj.UUID)
	w.Octree.RemoveUsing(*obj, obj.Node)
	w.Unlock()

	// update neighbors
	for _, neighbor := range neighbors {
		if neighbor.Type != types.ObjectTypePlayer && neighbor.Type != types.ObjectTypeNPC {
			continue
		}

		w.updateNeighbors(neighbor)
	}
}

func (w *World) hideObject(uuid string) {
	w.Lock()

	obj, ok := w.objects[uuid]
	if !ok {
		w.Unlock()
		return
	}

	neighbors := obj.Neighbors
	w.Octree.RemoveUsing(*obj, obj.Node)
	w.Unlock()

	for _, neighbor := range neighbors {
		if neighbor.Type != types.ObjectTypePlayer && neighbor.Type != types.ObjectTypeNPC {
			continue
		}

		w.updateNeighbors(neighbor)
	}
}

func (w *World) moveObjectTo(obj *types.GameObject) {
	w.RLock()
	w.Octree.RemoveUsing(*obj, obj.Node)
	obj.Node = w.Octree.Add(obj, types.Vector3f{obj.Position.X, obj.Position.Y, obj.Position.Z})
	w.RUnlock()
	w.updateNeighbors(obj)
}

func (w *World) getObjectsAt(position types.Vector3f) []*types.GameObject {
	w.Lock()
	defer w.Unlock()
	return w.Octree.ElementsAt(types.Vector3f(position))
}

func (w *World) updateNeighbors(obj *types.GameObject) {
	w.Lock()
	defer w.Unlock()

	radius := AREA_OF_INTEREST
	obj.Neighbors = nil
	center := obj.Position

	boxMin := types.Vector3f{center.X - radius, center.Y - radius, center.Z - radius}
	boxMax := types.Vector3f{center.X + radius, center.Y + radius, center.Z + radius}
	box := types.Box{Min: boxMin, Max: boxMax}

	elements := w.Octree.ElementsIn(box)

	for _, data := range elements {
		if obj.UUID != data.UUID {
			obj.Neighbors = append(obj.Neighbors, data)
		}

	}
}

func (w *World) updateNeighborsNearObject(obj *types.GameObject) {
	center := obj.Position

	radius := AREA_OF_INTEREST
	boxMin := types.Vector3f{center.X - radius, center.Y - radius, center.Z - radius}
	boxMax := types.Vector3f{center.X + radius, center.Y + radius, center.Z + radius}
	box := types.Box{Min: boxMin, Max: boxMax}

	w.Lock()
	elements := w.Octree.ElementsIn(box)
	w.Unlock()

	for _, data := range elements {
		w.updateNeighbors(data)
	}
}

func (w *World) getPlayersByPosition(position types.Vector3, radius float64) []*types.GameObject {
	w.RLock()
	defer w.RUnlock()

	neighbors := make([]*types.GameObject, 0)
	center := position

	boxMin := types.Vector3f{center.X - radius, center.Y - radius, center.Z - radius}
	boxMax := types.Vector3f{center.X + radius, center.Y + radius, center.Z + radius}
	box := types.Box{Min: boxMin, Max: boxMax}

	elements := w.Octree.ElementsIn(box)

	for _, data := range elements {
		if data.Type == types.ObjectTypePlayer {
			neighbors = append(neighbors, data)
		}
	}

	return neighbors
}

func (w *World) updateObjectVariation(obj *types.GameObject, variationIndex int32) {
	obj.VariationIndex = variationIndex
	w.updateNeighbors(obj)
	UpdateGameObjectVariationChannel <- &types.GameObjectVariation{Object: obj, VariationIndex: variationIndex}
}

func (w *World) broadcastSound(resource string, position types.Vector3, volume float32) {
	BroadcastSoundChannel <- &types.BroadcastSound{Resource: resource, Position: position, Volume: volume}
}

func (w *World) interactQueue(object *types.GameObject) {
	InteractQueue <- &types.InteractQueue{Object: object}
}

func (w *World) transformObjectRotation(object *types.GameObject, rotation types.Vector3) {
	object.Rotation = rotation
	UpdateTransformRotationChannel <- &types.TransformRotation{Object: object, Rotation: rotation}
}

func (w *World) dropItemOnGround(entity entity.Entity, position types.Vector3) {
	uuid, _ := uuid.NewUUID()

	destroyTime := time.Now().Add(time.Duration(10) * time.Second)
	object := &types.GameObject{
		Entity:      entity,
		Position:    position,
		UUID:        uuid.String(),
		Type:        types.ObjectTypeMapObject,
		DestroyTime: &destroyTime,
	}

	w.addObject(object)
	w.updateNeighbors(object)

	SpawnObjectChannel <- &types.SpawnObject{Object: object}
}

func (w *World) npcWalkTick() {
	for _, npc := range w.objects {
		if npc.Type != types.ObjectTypeNPC || npc.IsDead() {
			continue
		}

		// Destination set
		if npc.NextDestinationTime == nil || time.Now().After(*npc.NextDestinationTime) {
			if len(npc.Path) == 0 && len(npc.Waypoints) > 0 && npc.AttackTargetUUID == "" && !npc.IsReturningInProgress {
				waypoint := npc.GetNextRandomWaypoint()

				if waypoint != nil {

					w.RLock()
					npc.PathTargetAngleY = &waypoint[2]
					npc.SetDestination(waypoint[0], waypoint[1])
					w.RUnlock()

				}

			}
		}

		if len(npc.Path) > 0 {
			w.RLock()
			changed, finished := npc.MoveNPCWithWaypoints()
			w.RUnlock()

			if finished || changed {
				w.onWalkUpdates(npc)
			}

			if finished && npc.IsReturningInProgress {
				npc.IsReturningInProgress = false
			}

			if changed {
				UpdateMovementChannel <- npc
			}
		}

	}
}

func (w *World) npcRespawnTick() {
	for _, object := range w.objects {
		if object.Type != types.ObjectTypeNPC || object.NextSpawnTime == nil {
			continue
		}

		if time.Now().Before(*object.NextSpawnTime) {
			continue
		}

		W.Lock()
		object.Entity.Health = object.Entity.MaxHealth
		object.NextSpawnTime = nil
		object.Position = object.PositionSpawn
		object.Rotation = object.RotationSpawn
		object.SetNextTravelTime()
		W.Unlock()

		W.addObject(object)
		W.updateNeighbors(object)

		SpawnObjectChannel <- &types.SpawnObject{Object: object}
	}
}

func (w *World) npcResetCurrentAnimation(object *types.GameObject) {
	if object.CurrentAnimation == nil {
		return
	}

	UpdateAnimationChannel <- &types.Animation{Object: object, Name: *object.CurrentAnimation, IsStop: true}
	object.CurrentAnimation = nil
}

func (w *World) npcAttackTick() {
	for _, object := range w.objects {
		if object.Type != types.ObjectTypeNPC {
			continue
		}

		if !object.Entity.CanAgro || object.IsDead() || object.IsReturningInProgress {
			continue
		}

		attackRange := object.GetAttackRange()
		if attackRange == nil {
			continue
		}

		if object.AttackTargetUUID != "" {
			target, err := w.getObject(object.AttackTargetUUID)
			distanceFromSpawn := distance(types.Vector3{X: object.Waypoints[0][0], Y: object.Waypoints[0][2], Z: object.Waypoints[0][1]}, object.Position)

			if err != nil || target.IsDead() {
				fmt.Println("Lost target: ", distanceFromSpawn)

				object.ReleaseAttack()

				waypoint := object.GetNextRandomWaypoint()
				object.SetDestination(waypoint[0], waypoint[1])
				w.npcResetCurrentAnimation(object)
				continue
			}

			dist := distance(object.Position, target.Position)
			if dist <= *attackRange {

				object.Path = nil

				// Look at target
				if object.TargetPosition == nil || *object.TargetPosition != target.Position {
					lookAtRotation := object.LookAt(target)
					UpdateTransformRotationChannel <- &types.TransformRotation{Object: object, Rotation: lookAtRotation}
				}

				targetPosition := target.Position
				object.TargetPosition = &targetPosition

				// Skip attack frame
				if object.NextAttackTime != nil && time.Now().Before(*object.NextAttackTime) {
					continue
				}

				if object.IsReloadWeaponInProgress() {
					continue
				}

				if object.IsClipEmpty() {
					object.StartReloadWeapon()
					UpdateAnimationChannel <- &types.Animation{Object: object, Name: "Reloading", Speed: 1}
					continue
				}

				// Attack
				object.DecrementWeaponClip()
				object.AttackAttempts++
				attackSpeed := object.GetAttackSpeed()

				nextAttackTime := time.Now().Add(time.Duration(*attackSpeed*1000) * time.Millisecond)
				object.NextAttackTime = &nextAttackTime

				animation := object.GetInteractAnimation()
				if animation != "" {
					object.CurrentAnimation = &animation
					UpdateAnimationChannel <- &types.Animation{Object: object, Name: animation, Speed: 1}
				}

				go w.damageWithDelay(object, target, time.Duration(*attackSpeed*200))

				continue
			}

			w.npcResetCurrentAnimation(object)

			if len(object.Path) == 0 {
				fmt.Printf("NPC %s dist to target %f\n", object.Name, dist)
				isTargetNotReached := dist > float64(object.Entity.AttackRange)
				// TODO: check if NPC out of range from spawn
				if isTargetNotReached {
					object.SetDestination(target.Position.X, target.Position.Z)

					distanceFromSpawn := distance(*object.GetSpawnPoint(), object.Position)
					if distanceFromSpawn > 20 {
						object.Path = nil
						object.IsReturningInProgress = true
						object.ReleaseAttack()
						object.SetDestination(object.Waypoints[0][0], object.Waypoints[0][1])
						continue
					}
				}
			}
		}

		if object.AttackTargetUUID == "" {
			player, dist := w.findClosestPlayer(object)

			if player != nil && dist < 10 {
				object.AttackTargetUUID = player.UUID
				object.AttackAttempts = 0
				//object.Entity.Speed = 4
			}
		}
	}
}

func (w *World) damageWithDelay(source *types.GameObject, target *types.GameObject, delay time.Duration) {
	time.Sleep(delay)

	// random damage from 10 to 50
	isCrit := false
	damage := source.GetAttackMaxDamage()
	if damage == nil {
		return
	}

	target.TakeDamage(*damage)

	DamageChannel <- &types.Damage{Object: target, Amount: *damage, IsCrit: isCrit, HealthCurrent: int32(target.Entity.Health), HealthMax: int32(target.Entity.MaxHealth)}

	if target.IsDead() {
		source.AttackTargetUUID = ""
		source.TargetPosition = nil
		w.npcResetCurrentAnimation(source)

		go w.playerIsDead(target)
	}
}

func (w *World) playerIsDead(object *types.GameObject) {
	time.Sleep(4 * time.Second)

	teleport := w.getTeleport("main")
	object.Position = teleport.Position
	object.Health = object.MaxHealth
	TeleportObjectChannel <- &types.TeleportObject{Object: object, Position: teleport.Position, Rotation: teleport.Rotation}
}

func (w *World) getTeleport(name string) *LevelTeleport {
	for _, teleport := range w.teleports {
		if teleport.Name == name {
			return teleport
		}
	}
	return nil
}

func (w *World) mapObjectVariationTick() {
	for _, obj := range w.objects {
		if obj.NextVariation != nil && time.Now().After(obj.NextVariation.Time) {
			if obj.NextVariation.ResetHealth {
				obj.Entity.Health = obj.Entity.MaxHealth
			}
			w.updateObjectVariation(obj, obj.NextVariation.VariationIndex)
			obj.NextVariation = nil
		}
	}
}

func (w *World) mapObjectDestroyTick() {
	for _, obj := range w.objects {
		if obj.DestroyTime != nil && time.Now().After(*obj.DestroyTime) {
			w.removeObject(obj.UUID)
			DestroyObjectChannel <- &types.DestroyObject{Object: obj}
		}
	}
}

func (w *World) findClosestMapObjectByKind(gameObject *types.GameObject, kind types.ObjectKind) (*types.GameObject, float64) {
	w.Lock()
	defer w.Unlock()

	var closest *types.GameObject
	minDistance := math.MaxFloat64

	for _, point := range gameObject.Neighbors {
		if point.Kind == kind {
			dist := distance(gameObject.Position, point.Position)
			if dist < minDistance {
				minDistance = dist
				closest = point
			}
		}
	}

	return closest, minDistance
}

func (w *World) findClosestPlayer(gameObject *types.GameObject) (*types.GameObject, float64) {
	w.Lock()
	defer w.Unlock()

	var closest *types.GameObject
	minDistance := math.MaxFloat64

	for _, point := range gameObject.Neighbors {
		if point.Type != types.ObjectTypePlayer || point.Entity.Health <= 0 {
			continue
		}

		dist := distance(gameObject.Position, point.Position)
		if dist < minDistance {
			minDistance = dist
			closest = point
		}

	}

	return closest, minDistance
}

func (w *World) findLookedAtObjects(source *types.GameObject, targetName string, maxDistance, fieldOfViewAngle float64) []LookedAtObject {
	lookedAtObjectData := make([]LookedAtObject, 0)

	direction := types.Vector3{
		X: math.Cos(source.Rotation.Y) * math.Cos(source.Rotation.X),
		Y: math.Sin(source.Rotation.X),
		Z: math.Sin(source.Rotation.Y) * math.Cos(source.Rotation.X),
	}

	for _, point := range source.Neighbors {
		if point.InternalName == targetName {

			vectorToObject := types.Vector3{
				X: point.Position.X - source.Position.X,
				Y: point.Position.Y - source.Position.Y,
				Z: point.Position.Z - source.Position.Z,
			}

			length := math.Sqrt(math.Pow(vectorToObject.X, 2) + math.Pow(vectorToObject.Y, 2) + math.Pow(vectorToObject.Z, 2))
			vectorToObject.X /= length
			vectorToObject.Y /= length
			vectorToObject.Z /= length

			if length > maxDistance {
				continue
			}

			angle := math.Acos(direction.X*vectorToObject.X+direction.Y*vectorToObject.Y+direction.Z*vectorToObject.Z) * (180 / math.Pi)

			if angle <= fieldOfViewAngle/2 {
				lookedAtObjectData = append(lookedAtObjectData, LookedAtObject{GameObject: point, Distance: length})
			}
		}
	}

	sort.Slice(lookedAtObjectData, func(i, j int) bool {
		return lookedAtObjectData[i].Distance < lookedAtObjectData[j].Distance
	})

	return lookedAtObjectData
}

func distance(p1, p2 types.Vector3) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
}
