package structs

type alternative struct {
	path          path
	pressure      int64
	timeRemaining int64
	i             *valve
	j             *valve
}

type trip [][]*alternative
