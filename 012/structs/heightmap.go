package structs

import (
	"fmt"
	"math"
	"sort"
	"sync"
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
	return hm.getFewestStepsBetween(hm.start, hm.end)
}

func (hm *Heightmap) GetFewestStepsFromMinHeightToFinish() int {
	var wg sync.WaitGroup
	stepCounts := make([]int, 0, hm.xSize*hm.ySize)
	distance := make(chan int)

	for _, n := range hm.nodes {
		if n.height == 0 {
			wg.Add(1)
			go func(wg *sync.WaitGroup, out chan<- int, n *node) {
				defer wg.Done()
				out <- hm.getFewestStepsBetween(n, hm.end)
			}(&wg, distance, n)
		}
	}

	go func(wg *sync.WaitGroup, out chan<- int) {
		wg.Wait()
		close(out)
	}(&wg, distance)

	for dist := range distance {
		stepCounts = append(stepCounts, dist)
	}

	sort.Ints(stepCounts)

	return stepCounts[0]
}

func (hm *Heightmap) getFewestStepsBetween(from *node, to *node) int {
	visited := map[*node]int{}
	queue := queue{
		items: make([]queueNode, 0, hm.xSize*hm.ySize),
	}

	for _, n := range hm.nodes {
		visited[n] = math.MaxInt
	}

	visited[from] = 0
	queue.enqueue(from, 0)

	for queue.size() > 0 {
		current, distance := queue.dequeue()

		if current == to {
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

	return visited[to]
}
