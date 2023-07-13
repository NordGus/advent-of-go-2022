package structs

import "sync"

type Volcano struct {
	valves map[valveName]*valve
	start  *valve
}

func NewVolcano() Volcano {
	return Volcano{
		valves: make(map[valveName]*valve, 100),
		start:  nil,
	}
}

func (v *Volcano) AddValve(name string, flowRate int64, neighbors []string) {
	isStart := len(v.valves) == 0
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

	if isStart {
		v.start = val
	}
}

func (v *Volcano) ReleaseTheMostPressureWithin(timeLimit int64) int64 {
	wg := new(sync.WaitGroup)

	// fill all valves' shortest paths
	wg.Add(1)
	go v.fillValvesShortestPath(wg)
	wg.Wait()

	table := make(map[*valve]map[*valve]bool, len(v.valves))

	for _, from := range v.valves {
		table[from] = make(map[*valve]bool, len(v.valves))

		from.shortestPath.Print()

		for _, to := range v.valves {
			table[from][to] = true
		}
	}

	// for remaining > 0 {
	// 	segment := v.nextMoves(current, pressure, remaining)

	// 	if len(segment) == 0 {
	// 		break
	// 	}

	// 	path = append(path, segment...)
	// 	last := path[len(path)-1]

	// 	current = v.valves[last.to]
	// 	remaining = last.visitedAt
	// 	pressure = last.pressure

	// 	current.open()
	// }

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
	var startTime int64 = 0
	var pressure int64 = 0

	visited := make(map[valveName]bool, len(v.valves))
	moveQueue := make([]step, 0, len(v.valves))
	out := make(path, len(v.valves))

	if start.flowRate > 0 {
		startTime++
		pressure += start.flowRate
	}

	for name := range v.valves {
		out[name] = &step{
			from:      "-",
			to:        name,
			pressure:  pressure,
			visitedAt: startTime,
			openValve: name == start.name && start.flowRate > 0,
		}
	}

	visited[start.name] = true
	moveQueue = append(moveQueue, *out[start.name])

	for len(moveQueue) > 0 {
		current := moveQueue[0].to
		elapsed := moveQueue[0].visitedAt
		moveQueue = moveQueue[1:]

		for _, neighbor := range v.valves[current].neighbors {
			if visited[neighbor.name] {
				continue
			}

			visited[neighbor.name] = true

			move := out[neighbor.name]

			move.from = current
			move.visitedAt = elapsed + 1

			if neighbor.flowRate > 0 {
				move.pressure += neighbor.flowRate
				move.openValve = true
				move.visitedAt++
			}

			moveQueue = append(moveQueue, *move)
		}
	}

	return out
}
