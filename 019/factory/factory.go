package factory

import (
	"github.com/NordGus/advent-of-go-2022/019/shared/queue"
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
	top       tick
}

func NewFactory(blueprint Blueprint) Factory {
	return Factory{
		blueprint: blueprint,
		top: tick{
			state: state{OreRobots: 1},
			time:  0,
		},
	}
}

func (f *Factory) BlueprintID() int {
	return f.blueprint.id
}

func (f *Factory) QualityScoreDuring(duration int) int {
	var (
		states  = queue.New[tick]()
		visited = make(map[state]bool, 1000)
		robots  = []Resource{Geode, Obsidian, Clay, Ore}
	)

	_ = states.Enqueue(f.top, duration-f.top.time)

	for i := uint32(0); !states.IsEmpty(); i++ {
		current, _ := states.Pop()

		if current.state.Geode > f.top.state.Geode {
			f.top = current
		}

		if visited[current.state] || current.time == duration {
			continue
		}

		for i := 0; i < len(robots); i++ {
			if f.blueprint.robots[robots[i]].canBeBuilt(current.state.Ore, current.state.Clay, current.state.Obsidian) {
				build := tick{state: nextState(current.state, f.blueprint.robots[robots[i]]), time: current.time + 1}

				if !visited[build.state] {
					_ = states.Enqueue(build, duration-build.time)
				}
			}
		}

		produce := tick{state: nextState(current.state, robot{}), time: current.time + 1}

		if !visited[produce.state] {
			_ = states.Enqueue(produce, duration-produce.time)
		}

		visited[current.state] = true

		if i%1000 == 0 {
			runtime.GC()
		}
	}

	visited = map[state]bool{}
	states = nil

	runtime.GC()

	return f.top.state.Geode
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
