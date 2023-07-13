package structs

import (
	"fmt"
	"sort"
)

// step is simple struct to contain problem's state while walking trough the Volcano graph
type step struct {
	from      valveName
	to        valveName
	pressure  int64
	visitedAt int64
	openValve bool
}

// sortByTime is simple sorting interface to sort a slice representation of a path trough the Volcano graph by the moment each vertex is visited
type sortByTime []step

func (a sortByTime) Len() int           { return len(a) }
func (a sortByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByTime) Less(i, j int) bool { return a[i].visitedAt < a[j].visitedAt }

type path map[valveName]*step

func (p path) pathTo(destination *valve) path {
	out := make(path, len(p))

	return out
}

func (p path) Print() {
	steps := make([]step, 0, len(p))
	for _, s := range p {
		steps = append(steps, *s)
	}

	sort.Sort(sortByTime(steps))

	fmt.Println(steps)
}
