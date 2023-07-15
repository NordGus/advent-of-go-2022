package structs

import (
	"math"
)

type state struct {
	timeRemaining int64
	valve         valveName
	mask          uint32
}

// Volcano is a simple graph representation of problem's map layout
type Volcano struct {
	valves   map[valveName]*valve
	start    *valve
	cache    map[state]*float64
	openable []*valve
}

func NewVolcano() Volcano {
	return Volcano{
		valves: make(map[valveName]*valve, 100),
		start:  nil,
	}
}

func (v *Volcano) ParseValve(name string, flowRate int64, neighbors []string) {
	vName := valveName(name)
	val := v.valves[vName]

	if val == nil {
		val = &valve{
			name:    vName,
			tunnels: make(map[valveName]tunnel, defaultNeighborsMapSize),
		}

		v.valves[vName] = val
	}

	val.flowRate = flowRate

	for i := 0; i < len(neighbors); i++ {
		neighborValveName := valveName(neighbors[i])
		neighbor := v.valves[neighborValveName]

		if neighbor == nil {
			neighbor = &valve{
				name:    neighborValveName,
				tunnels: make(map[valveName]tunnel, defaultNeighborsMapSize),
			}

			v.valves[neighborValveName] = neighbor
		}

		val.tunnels[neighborValveName] = tunnel{to: neighbor, travelTime: 1}
	}

	if val.name == "AA" {
		v.start = val
	}
}

func (v *Volcano) ReleaseTheMostPressureWithin(timeLimit int64) int64 {
	v.cache = make(map[state]*float64, len(v.valves)*len(v.valves))

	out := v.exploreDepthFirstSearchMaxPressure(timeLimit, v.start, 0)

	return int64(out)
}

func (v *Volcano) exploreDepthFirstSearchMaxPressure(timeRemaining int64, source *valve, mask uint32) float64 {
	state := state{timeRemaining: timeRemaining, valve: source.name, mask: mask}

	if v.cache[state] != nil {
		return *v.cache[state]
	}

	var maxPressure float64 = 0

	for name, val := range v.valves {
		if name == source.name {
			continue
		}

		bit := uint32(1) << val.index

		if mask&bit != 0 {
			continue
		}

		remaining := timeRemaining - source.shortestPath.distanceTo(name) - 1

		if remaining <= 0 {
			continue
		}

		maxPressure = math.Max(maxPressure, v.exploreDepthFirstSearchMaxPressure(remaining, val, mask|bit)+float64(val.flowRate*remaining))
	}

	v.cache[state] = &maxPressure
	return maxPressure
}

func (v *Volcano) exploreBreathFirstSearch(start valveName) path {
	visited := make(map[valveName]bool, len(v.valves))
	moveQueue := make([]step, 0, len(v.valves))
	path := make(path, len(v.valves))

	visited[start] = true
	path[start] = step{from: "-", to: start, distance: 0}
	moveQueue = append(moveQueue, path[start])

	for len(moveQueue) > 0 {
		current := moveQueue[0]
		moveQueue = moveQueue[1:]

		for _, neighbor := range v.valves[current.to].tunnels {
			if visited[neighbor.to.name] {
				continue
			}

			visited[neighbor.to.name] = true

			move := step{
				from:      current.to,
				to:        neighbor.to.name,
				visitedAt: current.visitedAt + 1,
				distance:  current.distance + neighbor.travelTime,
			}

			moveQueue = append(moveQueue, move)
			path[neighbor.to.name] = move
		}
	}

	return path
}

func (v *Volcano) Simplify() *Volcano {
	volcano := &Volcano{
		valves:   make(map[valveName]*valve, len(v.valves)),
		openable: make([]*valve, 0, len(v.valves)-1),
	}

	for name, val := range v.valves {
		if val.flowRate > 0 || name == "AA" {
			volcano.valves[name] = &valve{
				name:     name,
				flowRate: val.flowRate,
				tunnels:  make(map[valveName]tunnel, len(v.valves)),
			}

		}

		if name == "AA" {
			volcano.start = volcano.valves[name]
		}
	}

	var index uint32 = 0

	for name, val := range volcano.valves {
		paths := v.exploreBreathFirstSearch(name)

		for neighborName, neighbor := range volcano.valves {
			if name == neighborName {
				continue
			}

			steps := paths.pathTo(neighborName)

			for i := 0; i < len(steps); i++ {
				current := steps[i].to

				if volcano.valves[current] != nil && current != neighborName {
					val.tunnels[current] = tunnel{
						to:         volcano.valves[current],
						travelTime: int64(i + 1),
					}
					break
				}

				if current == neighborName {
					val.tunnels[neighborName] = tunnel{
						to:         neighbor,
						travelTime: int64(i + 1),
					}
				}
			}
		}

		if name != "AA" {
			val.index = index
			volcano.openable = append(volcano.openable, val)
			index++
		}
	}

	for name, val := range volcano.valves {
		val.shortestPath = volcano.exploreBreathFirstSearch(name)
	}

	return volcano
}
