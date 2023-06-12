package structs

type MovementDirection int

const (
	RightMovement MovementDirection = iota
	LeftMovement
	UpMovement
	DownMovement
)

type Grid struct {
	head *head
	tail *tail
}

type Movement struct {
	Direction MovementDirection
	Amount    int
}

func NewGrid() Grid {
	h := head{
		x:       0,
		y:       0,
		history: make([][2]int, 0),
	}

	t := tail{
		x:       0,
		y:       0,
		history: make([][2]int, 0),
	}

	return Grid{
		head: &h,
		tail: &t,
	}
}

func (g *Grid) Move(move Movement) {
	for i := 0; i < move.Amount; i++ {
		switch move.Direction {
		case RightMovement:
			g.head.moveRight()
		case LeftMovement:
			g.head.moveLeft()
		case UpMovement:
			g.head.moveUp()
		case DownMovement:
			g.head.moveDown()
		default:
			panic("grid: unsupported movement direction")
		}

		g.tail.move(g.head)
	}
}

/*
func (g *Grid) Print() {
	fmt.Println(g.head)
	fmt.Println(g.tail)
}
*/

func (g *Grid) CountTailUniqueLocations() int {
	pos := map[[2]int]int{}

	for _, position := range g.tail.history {
		pos[position] += 1
	}

	return len(pos)
}
