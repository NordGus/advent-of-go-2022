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
	state   state
	skipped Resource
	time    int
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
		top     = tick{state: state{OreRobots: 1}, time: 0}
	)

	states = append(states, top)

	for i := uint32(0); len(states) > 0; i++ {
		var (
			built   = false
			current = states[0]
		)

		states = states[1:]

		if current.state.Geode > top.state.Geode {
			top = current
		}

		if visited[current.state] || current.time == duration {
			continue
		}

		for i := 0; i < len(robots); i++ {
			if current.time+1 == duration {
				break
			}

			if current.skipped == robots[i] {
				produce := tick{state: nextState(current.state, robot{}), time: current.time + 1}

				if !visited[produce.state] && produce.state.Geode >= top.state.Geode-1 {
					states = append(states, produce)
				}

				continue
			}

			if f.blueprint.robots[robots[i]].canBeBuilt(current.state.Ore, current.state.Clay, current.state.Obsidian) {
				build := tick{state: nextState(current.state, f.blueprint.robots[robots[i]]), time: current.time + 1}

				if !visited[build.state] && build.state.Geode >= top.state.Geode-1 {
					states = append(states, build)
				}

				built = true

				if robots[i] == Geode {
					break
				}

				noBuild := tick{state: nextState(current.state, robot{}), skipped: robots[i], time: current.time + 1}

				if !visited[noBuild.state] && noBuild.state.Geode >= top.state.Geode-1 {
					states = append(states, noBuild)
				}
			}
		}

		if !built {
			produce := tick{state: nextState(current.state, robot{}), time: current.time + 1}

			if !visited[produce.state] && produce.state.Geode >= top.state.Geode-1 {
				states = append(states, produce)
			}
		}

		visited[current.state] = true

		if i%10000 == 0 {
			runtime.GC()
		}
	}

	runtime.GC()

	return top.state.Geode
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
