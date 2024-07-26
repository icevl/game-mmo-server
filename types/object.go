package types

import (
	"math"
	"math/rand"
	"server/entity"
	pbglobal "server/proto"
	"server/proto/actionpb"
	"server/proto/transformpb"
	"server/utils"
	"time"
)

type ObjectType = string
type ObjectKind int

const (
	ObjectTypePlayer           ObjectType = "player"
	ObjectTypeNPC              ObjectType = "npc"
	ObjectTypeVariantMapObject ObjectType = "variant_map_object" // tree, rock, etc. placed in server side
	ObjectTypeMapObject        ObjectType = "loot_object"        // loot object placed in client side
)

const (
	ObjectKindDefault  = 0
	ObjectKindNPC      = 1
	ObjectKindTeleport = 2
	ObjectKindTree     = 3
	ObjectKindOre      = 4
)

type EventPayloadOptions struct {
	IsSelf bool
}

type GameObjectVariation struct {
	Object         *GameObject
	VariationIndex int32
}

type NextVariation struct {
	VariationIndex int32
	Time           time.Time
	ResetHealth    bool
}

type GameObject struct {
	entity.Entity

	UUID string
	Node *Node

	Kind          ObjectKind
	Position      Vector3
	Rotation      Vector3
	PositionSpawn Vector3
	RotationSpawn Vector3

	Type ObjectType

	VariationIndex int32

	Neighbors        []*GameObject
	Waypoints        [][3]float64 // X, Z, RotationY
	Path             [][3]float64 // Path of last waypoint
	PathTargetAngleY *float64     // Target waypoint Rotation Y
	TargetPosition   *Vector3     // NPC attack target

	CurrentAnimation        *string
	AttackTargetUUID        string
	AttackAttempts          int32
	IsReturningInProgress   bool       // Return to spawn point for NPC
	NextDestinationTime     *time.Time // Next time for destination set
	NextAttackTime          *time.Time // Next time for attack set
	NextStepTime            *time.Time // Timer for next step
	NextTransformUpdateTime *time.Time
	NextSpawnTime           *time.Time
	NextVariation           *NextVariation
	DestroyTime             *time.Time // Time to destroy object (loot, etc.)
}

func (o *GameObject) SetNextTravelTime() {
	o.NextStepTime = nil
	NextDestinationTime := time.Now().Add(time.Second * time.Duration(rand.Intn(120)+15))
	o.NextDestinationTime = &NextDestinationTime
}

func (o *GameObject) SetDestination(x, z float64) {
	path, err := utils.GetPath(o.Position.X, o.Position.Z, x, z)
	//fmt.Println("Path: ", path)
	if err != nil {
		//fmt.Printf("Path not found for x: %f z: %f\n", x, z)
		return
	}

	o.Path = path
}

func (o *GameObject) GetPlayersNearby() []*GameObject {
	players := make([]*GameObject, 0)

	for _, neighbor := range o.Neighbors {
		if neighbor.Type == ObjectTypePlayer {
			players = append(players, neighbor)
		}
	}

	return players
}

func (o *GameObject) getNPCsNearby() []*GameObject {
	npcs := make([]*GameObject, 0)

	for _, neighbor := range o.Neighbors {
		if neighbor.Type == ObjectTypeNPC {
			npcs = append(npcs, neighbor)
		}
	}

	return npcs
}

func (o *GameObject) GetMapObjectsNearby() []*GameObject {
	objects := make([]*GameObject, 0)

	for _, neighbor := range o.Neighbors {
		if neighbor.Type == ObjectTypeVariantMapObject {
			objects = append(objects, neighbor)
		}
	}

	return objects
}

func (o *GameObject) ScheduleRespawn() {
	if o.Entity.RespawnInterval == 0 {
		return
	}

	// NPC Respawn
	if o.Type == ObjectTypeNPC {
		nextSpawnTime := time.Now().Add(time.Duration(o.Entity.RespawnInterval) * time.Second)
		o.NextSpawnTime = &nextSpawnTime
		return
	}

	// Map object respawn
	o.NextVariation = &NextVariation{
		VariationIndex: 0,
		Time:           time.Now().Add(time.Duration(o.Entity.RespawnInterval) * time.Second),
		ResetHealth:    true,
	}
}

func (o *GameObject) TakeDamage(amount int32) {
	o.Entity.Health -= amount
	if o.Entity.Health <= 0 {
		o.Entity.Health = 0

		// Release current agro, path and respawn
		if o.Type == ObjectTypeNPC {
			o.ScheduleRespawn()
			o.TargetPosition = nil
			o.AttackTargetUUID = ""
			o.Path = nil
		}

	}

}

func (o *GameObject) GetAttackRange() *float64 {

	// Right hand
	if o.Entity.EquippedItems.RightHand.AttackRange > 0 {
		attackRange := float64(o.Entity.EquippedItems.RightHand.AttackRange)
		return &attackRange
	}

	// None Humanoid NPC
	if o.Entity.AttackRange > 0 {
		attackRange := float64(o.Entity.AttackRange)
		return &attackRange
	}

	return nil
}

func (o *GameObject) GetAttackMaxDamage() *int32 {

	// Right hand
	if o.Entity.EquippedItems.RightHand.AttackDamage > 0 {
		attackMaxDamage := int32(o.Entity.EquippedItems.RightHand.AttackDamage)
		return &attackMaxDamage
	}

	// None Humanoid NPC
	if o.Entity.AttackDamage > 0 {
		attackMaxDamage := int32(o.Entity.AttackDamage)
		return &attackMaxDamage
	}

	return nil
}

func (o *GameObject) GetAttackSpeed() *float64 {
	if o.Entity.EquippedItems.RightHand.AttackSpeed > 0 {
		attackSpeed := float64(o.Entity.EquippedItems.RightHand.AttackSpeed)
		return &attackSpeed
	}

	if o.Entity.AttackSpeed > 0 {
		attackSpeed := float64(o.Entity.AttackSpeed)
		return &attackSpeed
	}

	return nil
}

func (o *GameObject) ReleaseAttack() {
	o.AttackTargetUUID = ""
	o.TargetPosition = nil
}

func (o *GameObject) DecrementWeaponClip() {
	if o.Entity.EquippedItems.RightHand.ClipSize == 0 {
		return
	}

	o.Entity.EquippedItems.RightHand.Clip--
}

func (o *GameObject) IsClipEmpty() bool {
	if o.Entity.EquippedItems.RightHand.ClipSize == 0 {
		return false
	}

	return o.Entity.EquippedItems.RightHand.Clip == 0
}

func (o *GameObject) LookAt(target *GameObject) Vector3 {
	direction := Vector3{
		X: target.Position.X - o.Position.X,
		Y: 0,
		Z: target.Position.Z - o.Position.Z,
	}

	directionMagnitude := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y + direction.Z*direction.Z)
	normalizedDirection := Vector3{
		X: direction.X / directionMagnitude,
		Y: direction.Y / directionMagnitude,
		Z: direction.Z / directionMagnitude,
	}
	angle := math.Atan2(normalizedDirection.X, normalizedDirection.Z) * (180.0 / math.Pi)
	rotation := Vector3{X: 0, Y: angle, Z: 0}

	o.Rotation = rotation

	return rotation
}

func (o *GameObject) IsDead() bool {
	return o.Entity.Health <= 0
}

func (o *GameObject) GetNextRandomWaypoint() *[3]float64 {
	if len(o.Waypoints) == 0 {
		return nil
	}

	selectedWaypoints := make([]*[3]float64, 0)
	for i := range o.Waypoints {
		waypoint := &o.Waypoints[i]
		if math.Floor(waypoint[0]) == math.Floor(o.Position.X) && math.Floor(waypoint[1]) == math.Floor(o.Position.Z) {
			continue
		}
		selectedWaypoints = append(selectedWaypoints, waypoint)
	}

	if len(selectedWaypoints) == 0 {
		return nil
	}

	waypointIndex := rand.Intn(len(selectedWaypoints))
	point := selectedWaypoints[waypointIndex]

	return &[3]float64{math.Floor(point[0]), math.Floor(point[1]), point[2]}
}

func (o *GameObject) GetSpawnPoint() *Vector3 {
	if len(o.Waypoints) == 0 {
		return nil
	}

	return &Vector3{X: o.Waypoints[0][0], Y: o.Waypoints[0][2], Z: o.Waypoints[0][1]}
}

func (o *GameObject) MoveNPCWithWaypoints() (bool, bool) {
	if len(o.Path) == 0 {
		return false, false
	}

	if o.NextStepTime != nil {
		if time.Now().Before(*o.NextStepTime) {
			return false, false
		}
	}

	node := o.Path[0]

	directionX := node[0] - o.Position.X
	directionZ := node[2] - o.Position.Z
	angleY := math.Atan2(directionX, directionZ) * 180 / math.Pi

	o.Rotation.Y = angleY

	o.Position.X = node[0]
	o.Position.Y = node[1]
	o.Position.Z = node[2]

	o.Path = o.Path[1:]

	if len(o.Path) > 0 {
		nextNode := o.Path[0]
		distance := math.Sqrt(math.Pow(float64(nextNode[0]-node[0]), 2) + math.Pow(float64(nextNode[2]-node[2]), 2))
		sleepTime := int64(distance/float64(o.Speed)*1000.0) - 20 // 20ms for processing
		nextStep := time.Now().Add(time.Duration(sleepTime) * time.Millisecond)
		o.NextStepTime = &nextStep

		return true, false
	} else {
		o.SetNextTravelTime()
		return true, true
	}
}

func (o *GameObject) getNetworkTransformPayload() *actionpb.Action {
	msg := &transformpb.Transform{
		UUID:     o.UUID,
		Speed:    o.Speed,
		Position: &pbglobal.Vector3M{X: float32(o.Position.X), Y: float32(o.Position.Y), Z: float32(o.Position.Z)},
		Rotation: &pbglobal.Vector3M{X: float32(o.Rotation.X), Y: float32(o.Rotation.Y), Z: float32(o.Rotation.Z)},
	}

	return &actionpb.Action{
		Action: &actionpb.Action_Transform{
			Transform: msg,
		},
	}
}
