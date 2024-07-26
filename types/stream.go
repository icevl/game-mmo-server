package types

type TransformRotation struct {
	Object   *GameObject
	Rotation Vector3
}

type TeleportObject struct {
	Object   *GameObject
	Position Vector3
	Rotation Vector3
}

type Animation struct {
	Object *GameObject
	Name   string
	Speed  float32
	IsStop bool
}

type StopAnimation struct {
	Object *GameObject
	Name   string
}

type Damage struct {
	Object        *GameObject
	Amount        int32
	HealthCurrent int32
	HealthMax     int32
	IsCrit        bool
}

type BroadcastSound struct {
	Resource string
	Position Vector3
	Volume   float32
}

type DestroyObject struct {
	Object *GameObject
}

type SpawnObject struct {
	Object *GameObject
}

type InteractQueue struct {
	Object *GameObject
}
