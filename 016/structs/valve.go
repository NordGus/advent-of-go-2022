package structs

const (
	defaultNeighborsMapSize = 5
)

// valveName is a simple alias for the string type to make some typing more readable
type valveName string

// valve is a vertex for the Volcano graph representing a room from the problem
type valve struct {
	name         valveName
	flowRate     int64
	neighbors    map[valveName]*valve
	shortestPath path
	index        int
}
