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
	tunnels      map[valveName]tunnel
	shortestPath path
	index        uint32
}

type tunnel struct {
	to         *valve
	travelTime int64
}
