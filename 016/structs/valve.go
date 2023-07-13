package structs

const (
	defaultNeighborsMapSize = 5
)

type valveName string

type valve struct {
	name         valveName
	flowRate     int64
	neighbors    map[valveName]*valve
	shortestPath path
}
