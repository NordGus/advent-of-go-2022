package structs

import (
	"fmt"
	"math"
)

type nodeRole int

const (
	StartNode nodeRole = iota
	EndNode
	NormalNode
)

type Heightmap struct {
	xSize int
	ySize int
	start *node
	end   *node
	nodes map[[2]int]*node
}

func NewHeightmap(xSize int, ySize int) Heightmap {
	return Heightmap{
		xSize: xSize,
		ySize: ySize,
		nodes: map[[2]int]*node{},
	}
}

func (hm *Heightmap) AddNode(height int, x int, y int, role nodeRole) error {
	pos := [2]int{x, y}

	if hm.nodes[pos] != nil {
		return fmt.Errorf("node [%v,%v] already exists", x, y)
	}

	n := node{height: height, x: x, y: y, edges: make([]*node, 0)}

	if role == StartNode {
		hm.start = &n
	}

	if role == EndNode {
		hm.end = &n
	}

	hm.nodes[pos] = &n

	return nil
}

func (hm *Heightmap) BuildGraph() {
	for i := 0; i < hm.xSize; i++ {
		for j := 0; j < hm.ySize; j++ {
			node := hm.nodes[[2]int{i, j}]

			up := hm.nodes[[2]int{node.x - 1, node.y}]
			down := hm.nodes[[2]int{node.x + 1, node.y}]
			left := hm.nodes[[2]int{node.x, node.y - 1}]
			right := hm.nodes[[2]int{node.x, node.y + 1}]

			node.addEdge(up)
			node.addEdge(down)
			node.addEdge(left)
			node.addEdge(right)
		}
	}
}

func (hm *Heightmap) GetFewestStepsFromStartToFinish() int {
	visited := map[*node]int{}
	queue := queue{
		items: make([]queueNode, 0, 1_000),
	}

	for _, n := range hm.nodes {
		visited[n] = math.MaxInt
	}

	visited[hm.start] = 0
	queue.enqueue(hm.start, 0)

	for queue.size() > 0 {
		current, distance := queue.dequeue()

		if current == hm.end {
			queue.clear()
			break
		}

		for _, n := range current.edges {
			if visited[n] == math.MaxInt {
				visited[n] = distance + 1
				queue.enqueue(n, distance+1)
			}
		}
	}

	return visited[hm.end]
}
