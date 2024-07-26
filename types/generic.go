package types

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type Vector3f [3]float64

// Box Defines an axis aligned rectangular solid.
type Box struct {
	Min Vector3f
	Max Vector3f
}

type Node struct {
	box         Box
	point       *Vector3f
	elements    []*GameObject
	hasChildren bool
	children    []*Node
}
