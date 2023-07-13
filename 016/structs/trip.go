package structs

type alternative struct {
	path          []step
	pressure      int64
	timeRemaining int64
	openedValves  map[valveName]bool
	i             valve
	j             valve
}

type trip [][]*alternative
