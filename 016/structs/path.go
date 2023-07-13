package structs

import (
	"fmt"
	"sort"
)

type step struct {
	from      valveName
	to        valveName
	pressure  int64
	visitedAt int64
	openValve bool
}

type sortByTime []step

func (a sortByTime) Len() int           { return len(a) }
func (a sortByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByTime) Less(i, j int) bool { return a[i].visitedAt < a[j].visitedAt }

type path map[valveName]*step

func (p path) Print() {
	steps := make([]step, 0, len(p))
	for _, s := range p {
		steps = append(steps, *s)
	}

	sort.Sort(sortByTime(steps))

	fmt.Println(steps)
}
