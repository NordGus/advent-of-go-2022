package structs

import (
	"fmt"
	"sync"
)

// Volcano is a simple graph representation of problem's map layout
type Volcano struct {
	valves     map[valveName]*valve
	start      *valve
	valveIndex int
}

func NewVolcano() Volcano {
	return Volcano{
		valves: make(map[valveName]*valve, 100),
		start:  nil,
	}
}

func (v *Volcano) AddValve(name string, flowRate int64, neighbors []string) {
	vName := valveName(name)
	val := v.valves[vName]

	if val == nil {
		val = &valve{
			name:      vName,
			neighbors: make(map[valveName]*valve, defaultNeighborsMapSize),
		}

		v.valves[vName] = val
	}

	val.flowRate = flowRate
	val.index = v.valveIndex
	v.valveIndex++

	for i := 0; i < len(neighbors); i++ {
		neighborValveName := valveName(neighbors[i])
		neighbor := v.valves[neighborValveName]

		if neighbor == nil {
			neighbor = &valve{
				name:      neighborValveName,
				neighbors: make(map[valveName]*valve, defaultNeighborsMapSize),
			}

			v.valves[neighborValveName] = neighbor
		}

		val.neighbors[neighborValveName] = neighbor
	}

	if val.index == 0 {
		v.start = val
	}
}

func (v *Volcano) ReleaseTheMostPressureWithin(timeLimit int64) int64 {
	var out *room
	wg := new(sync.WaitGroup)
	trip := make(trip, len(v.valves))

	// fill all valves' shortest paths
	wg.Add(1)
	go v.fillValvesShortestPath(wg)
	wg.Wait()

	// build trip starting state
	for _, from := range v.valves {
		trip[from.index] = make([]*room, len(v.valves))

		for _, to := range v.valves {
			trip[from.index][to.index] = &room{
				valve:         *to,
				timeRemaining: timeLimit,
				openedValves:  make(map[valveName]bool, len(v.valves)),
			}
		}
	}

	// build solution
	for i := 0; i < len(v.valves); i++ {
		for j := 0; j < len(v.valves); j++ {
			var top *room
			var diagonal *room

			left := trip[i][0]

			if i == 0 && j == 0 {
				top = trip[i][j]
				diagonal = trip[i][j]
			}

			if i == 0 && j > 0 {
				top = trip[i][0]
				diagonal = trip[i][0]
			}

			if i > 0 && j == 0 {
				left = trip[j][i]
				diagonal = trip[i-1][j]
			}

			if top == nil {
				top = trip[i-1][j]
			}

			if diagonal == nil {
				diagonal = trip[i-1][j-1]
			}

			fromLeft := trip[i][j].travel(*left, v)
			fromTop := trip[i][j].travel(*top, v)
			fromDiagonal := trip[i][j].travel(*diagonal, v)

			if (fromLeft.pressureReleased * fromLeft.timeRemaining) > (trip[i][j].pressureReleased * trip[i][j].timeRemaining) {
				trip[i][j] = &fromLeft
			}

			if (fromTop.pressureReleased * fromTop.timeRemaining) > (trip[i][j].pressureReleased * trip[i][j].timeRemaining) {
				trip[i][j] = &fromTop
			}

			if (fromDiagonal.pressureReleased * fromDiagonal.timeRemaining) > (trip[i][j].pressureReleased * trip[i][j].timeRemaining) {
				trip[i][j] = &fromDiagonal
			}

			if i == 1 {
				fmt.Printf("current[%v][%v] %v\n", i, j, trip[i][j].valve.name)
				fmt.Printf("\t%v\n", trip[i][j].path)
				fmt.Printf("\tpressureReleased %v\n", trip[i][j].pressureReleased)
				fmt.Printf("\ttimeRemaining %v\n", trip[i][j].timeRemaining)
				fmt.Printf("\topenedValves %v\n", trip[i][j].openedValves)
			}

			if out == nil || trip[i][j].pressureReleased > out.pressureReleased {
				out = trip[i][j]
			}
		}
	}

	fmt.Println(out.pressureReleased)

	return 0
}

func (v *Volcano) fillValvesShortestPath(wg *sync.WaitGroup) {
	for _, val := range v.valves {
		wg.Add(1)
		go func(v *Volcano, val *valve, wg *sync.WaitGroup) {
			val.shortestPath = v.exploreFrom(val)
			wg.Done()
		}(v, val, wg)
	}

	wg.Done()
}

func (v *Volcano) exploreFrom(start *valve) path {
	visited := make(map[valveName]bool, len(v.valves))
	moveQueue := make([]step, 0, len(v.valves))
	path := make(path, len(v.valves))

	visited[start.name] = true
	path[start.name] = step{from: "-", to: start.name}
	moveQueue = append(moveQueue, path[start.name])

	for len(moveQueue) > 0 {
		current := moveQueue[0]
		moveQueue = moveQueue[1:]

		for _, neighbor := range v.valves[current.to].neighbors {
			if visited[neighbor.name] {
				continue
			}

			visited[neighbor.name] = true

			move := step{
				from:      current.to,
				to:        neighbor.name,
				visitedAt: current.visitedAt + 1,
			}

			moveQueue = append(moveQueue, move)
			path[neighbor.name] = move
		}
	}

	return path
}
