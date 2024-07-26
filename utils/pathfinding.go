package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"sort"
)

type GridData [][][2]float32

type Node struct {
	X, Z, Y float64
	Tag     int
}

type AStarNode struct {
	Node      *Node
	Cost      int
	Prev      *AStarNode
	Heuristic int
}

const FILE_PATH = "/Users/ice/MMO/Assets/Scenes/LE_ExampleGame_navgrid.json" //LE_ExampleGame_navgrid | LE_ExampleEditor_navgrid.json

const (
	offsetX = -1000
	offsetZ = -1000
)

const maxIterations = 1000

var grid [][]*Node

func init() {
	var err error
	grid, err = loadGridData(FILE_PATH)
	if err != nil {
		log.Fatal("Error int grid file reading: ", err)
	}
}

func loadGridData(filename string) ([][]*Node, error) {
	var gridData GridData

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &gridData)
	if err != nil {
		return nil, err
	}

	grid := make([][]*Node, len(gridData))

	for i := range gridData {
		grid[i] = make([]*Node, len(gridData[i]))
		for j := range gridData[i] {
			nodeData := gridData[i][j]
			grid[i][j] = &Node{X: float64(i), Z: float64(j), Y: float64(nodeData[1]), Tag: int(nodeData[0])}
		}
	}

	return grid, nil
}

func GetPath(startX, startZ, endX, endZ float64) ([][3]float64, error) {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return nil, errors.New("empty grid")
	}

	oX := int(offsetX) * -1
	oZ := int(offsetZ) * -1

	startNode := grid[int(startX)+oX][int(startZ)+oZ]
	endNode := grid[int(endX)+oX][int(endZ)+oZ]

	path, err := aStar(startNode, endNode, grid)

	if err != nil {
		return nil, err
	}

	result := make([][3]float64, 0)
	for _, node := range path {
		result = append(result, [3]float64{node.X + offsetX, node.Y, node.Z + offsetZ})
	}

	return result, nil
}

func aStar(start, goal *Node, grid [][]*Node) ([]*Node, error) {
	if !isWalkable(start) {
		return nil, errors.New("start node is not walkable")
	}

	if !isWalkable(goal) {
		return nil, errors.New("goal node is not walkable")
	}

	var closedSet []*AStarNode
	var openSet = []*AStarNode{{Node: start, Cost: 0, Heuristic: manhattanDistance(start, goal)}}
	iter := 0

	for len(openSet) > 0 {
		iter++
		if iter > maxIterations {
			return nil, errors.New("aStar: maximum iterations reached")
		}

		current := openSet[0]

		sort.Slice(openSet, func(i, j int) bool {
			return openSet[i].Heuristic < openSet[j].Heuristic
		})

		openSet = openSet[1:]

		if current.Node == goal {
			var path []*Node
			for current != nil {
				path = append([]*Node{current.Node}, path...)
				current = current.Prev
			}
			return path, nil
		}

		closedSet = append(closedSet, current)

		neighbors := getNeighbors(current.Node, grid)
		for _, neighbor := range neighbors {
			if containsAStarNode(closedSet, neighbor) {
				continue
			}

			tentativeCost := current.Cost + 1

			neighborAStar := findAStarNodeInSlice(openSet, neighbor)

			if neighborAStar == nil {
				neighborAStar = &AStarNode{Node: neighbor}
				openSet = append(openSet, neighborAStar)
			}

			if tentativeCost >= neighborAStar.Cost && neighborAStar.Cost != 0 {
				continue
			}

			neighborAStar.Prev = current
			neighborAStar.Cost = tentativeCost
			neighborAStar.Heuristic = tentativeCost + manhattanDistance(neighbor, goal)
		}
	}

	return nil, errors.New("aStar: path not found")
}

// проверяет, содержится ли искомый узел в передаваемом слайсе AStarNode
func containsAStarNode(nodes []*AStarNode, targetNode *Node) bool {
	for _, node := range nodes {
		if node.Node == targetNode {
			return true
		}
	}
	return false
}

func findAStarNodeInSlice(nodes []*AStarNode, targetNode *Node) *AStarNode {
	for _, node := range nodes {
		if node.Node == targetNode {
			return node
		}
	}

	return nil
}

func isWalkable(n *Node) bool {
	return n.Tag != 0
}

func manhattanDistance(nodeA, nodeB *Node) int {
	dx := math.Abs(nodeA.X - nodeB.X)
	dz := math.Abs(nodeA.Z - nodeB.Z)
	return int(dx + dz)
}

func getNeighbors(node *Node, grid [][]*Node) []*Node {
	var neighbors []*Node
	directions := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

	for _, d := range directions {
		y := int(node.X) + d[0]
		x := int(node.Z) + d[1]

		if y >= 0 && x >= 0 && y < len(grid) && x < len(grid[y]) {
			nextNode := grid[y][x]
			if isWalkable(nextNode) {
				neighbors = append(neighbors, nextNode)
			}
		}
	}
	return neighbors
}
