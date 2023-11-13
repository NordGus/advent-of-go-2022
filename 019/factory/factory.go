package factory

import (
	"runtime"
)

type state struct {
	OreRobots      int
	ClayRobots     int
	ObsidianRobots int
	GeodeRobots    int

	Ore      int
	Clay     int
	Obsidian int
	Geode    int
}

type tick struct {
	state state
	time  int
}

type Factory struct {
	blueprint Blueprint
}

func NewFactory(blueprint Blueprint) Factory {
	return Factory{
		blueprint: blueprint,
	}
}

func (f *Factory) BlueprintID() int {
	return f.blueprint.id
}

func (f *Factory) QualityScoreDuring(duration int) int {
	var (
		states  = make([]tick, 0, 100000)
		visited = make(map[state]bool, 100000)
		robots  = []Resource{Geode, Obsidian, Clay, Ore}
		top     = state{OreRobots: 1}
	)

	states = append(states, tick{state: top})

	for i := uint32(0); len(states) > 0; i++ {
		current := states[0]
		states = states[1:]

		if current.state.Geode > top.Geode {
			top = current.state
		}

		if visited[current.state] {
			continue
		}

		if current.time == duration {
			continue
		}

		for i := 0; i < len(robots); i++ {
			if f.blueprint.robots[robots[i]].canBeBuilt(current.state.Ore, current.state.Clay, current.state.Obsidian) {
				build := tick{state: nextState(current.state, f.blueprint.robots[robots[i]]), time: current.time + 1}

				if !visited[build.state] {
					states = append(states, build)
				}
			}
		}

		produce := tick{state: nextState(current.state, robot{}), time: current.time + 1}

		if !visited[produce.state] {
			states = append(states, produce)
		}

		visited[current.state] = true

		if i%10000 == 0 {
			runtime.GC()
		}
	}

	runtime.GC()

	return top.Geode
}

func nextState(prev state, rbt robot) state {
	next := prev

	if rbt.Resource != Invalid {
		next.Ore -= rbt.Ore
		next.Clay -= rbt.Clay
		next.Obsidian -= rbt.Obsidian
	}

	next.Ore += prev.OreRobots
	next.Clay += prev.ClayRobots
	next.Obsidian += prev.ObsidianRobots
	next.Geode += prev.GeodeRobots

	if rbt.Resource == Ore {
		next.OreRobots += 1
	}

	if rbt.Resource == Clay {
		next.ClayRobots += 1
	}

	if rbt.Resource == Obsidian {
		next.ObsidianRobots += 1
	}

	if rbt.Resource == Geode {
		next.GeodeRobots += 1
	}

	return next
}
