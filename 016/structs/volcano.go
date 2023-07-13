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
			index:     v.valveIndex,
		}

		v.valves[vName] = val
		v.valveIndex++
	}

	val.flowRate = flowRate

	for i := 0; i < len(neighbors); i++ {
		neighborValveName := valveName(neighbors[i])
		neighbor := v.valves[neighborValveName]

		if neighbor == nil {
			neighbor = &valve{
				name:      neighborValveName,
				neighbors: make(map[valveName]*valve, defaultNeighborsMapSize),
				index:     v.valveIndex,
			}

			v.valves[neighborValveName] = neighbor
			v.valveIndex++
		}

		val.neighbors[neighborValveName] = neighbor
	}

	if val.index == 0 {
		v.start = val
	}
}

func (v *Volcano) ReleaseTheMostPressureWithin(timeLimit int64) int64 {
	var out *alternative
	wg := new(sync.WaitGroup)
	trip := make(trip, len(v.valves))

	// fill all valves' shortest paths
	wg.Add(1)
	go v.fillValvesShortestPath(wg)
	wg.Wait()

	// build trip starting state
	for _, from := range v.valves {
		trip[from.index] = make([]*alternative, len(v.valves))

		for _, to := range v.valves {
			trip[from.index][to.index] = &alternative{
				i:             *from,
				j:             *to,
				timeRemaining: timeLimit,
			}
		}
	}

	// build solution
	for i := 0; i < len(v.valves); i++ {
		for j := 0; j < len(v.valves); j++ {
			current := trip[i][j]
			remaining := current.timeRemaining
			pressure := current.pressure
			opened := make(map[valveName]bool, len(v.valves))

			path := current.i.shortestPath.pathTo(current.j.name)[1:]

			for i := 0; i < len(path); i++ {
				val := v.valves[path[i].to]
				remaining--

				if val.flowRate > 0 && !opened[val.name] {
					remaining--
					opened[val.name] = true
					pressure += val.flowRate * remaining
				}
			}

			if i == 0 {
				current.path = path
				current.pressure = pressure
				current.timeRemaining = remaining
				current.openedValves = opened

				fmt.Println(current.i.name, current.j.name, current.path, current.pressure, current.timeRemaining, current.openedValves)
				continue
			}

			if j == 0 {
				continue
			}
		}

		if i == 0 {
			break
		}
	}

	fmt.Println(out)

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
