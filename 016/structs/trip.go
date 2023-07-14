package structs

type room struct {
	valve            valve
	path             []step
	pressureReleased int64
	timeRemaining    int64
	openedValves     map[valveName]bool
}

type trip [][]*room

func (r room) travel(prev room, volcano *Volcano) room {
	elapsed := r.timeRemaining - prev.timeRemaining
	remaining := prev.timeRemaining
	pressure := prev.pressureReleased
	path := prev.path
	opened := make(map[valveName]bool, len(prev.openedValves))

	for name, open := range prev.openedValves {
		opened[name] = open
	}

	path = append(path, prev.valve.shortestPath.pathTo(r.valve.name, elapsed)...)

	for i := len(prev.path); i < len(path); i++ {
		val := volcano.valves[path[i].to]
		remaining--

		if val.flowRate > 0 && !opened[val.name] {
			remaining--
			opened[val.name] = true
			pressure += val.flowRate * remaining
		}
	}

	return room{
		valve:            r.valve,
		path:             path,
		pressureReleased: pressure,
		timeRemaining:    remaining,
		openedValves:     opened,
	}
}
