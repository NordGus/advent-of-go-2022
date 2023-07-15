package structs

import (
	"sort"
)

// step is simple struct to contain problem's state while walking trough the Volcano graph
type step struct {
	from      valveName
	to        valveName
	visitedAt int64
	distance  int64
}

// sortByTime is simple sorting interface to sort a slice representation of a path trough the Volcano graph by the moment each vertex is visited
type sortByTime []step

func (a sortByTime) Len() int           { return len(a) }
func (a sortByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByTime) Less(i, j int) bool { return a[i].visitedAt < a[j].visitedAt }

type path map[valveName]step

func (p path) pathTo(destination valveName) []step {
	out := make([]step, 0, len(p))

	current := p[destination]

	for current.from != "-" {
		s := step{
			from:      current.from,
			to:        current.to,
			visitedAt: current.visitedAt,
		}

		out = append(out, s)

		current = p[current.from]
	}

	sort.Sort(sortByTime(out))

	return out
}

func (p path) distanceTo(destination valveName) int64 {
	return p[destination].distance
}
