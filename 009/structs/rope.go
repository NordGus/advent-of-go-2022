package structs

type MovementDirection int

const (
	RightMovement MovementDirection = iota
	LeftMovement
	UpMovement
	DownMovement

	knotCount = 10
)

type Rope struct {
	knots []*knot
}

type Movement struct {
	Direction MovementDirection
	Amount    int
}

func NewRope() Rope {
	knots := make([]*knot, knotCount)

	for i := 0; i < knotCount; i++ {
		knots[i] = &knot{
			x:       0,
			y:       0,
			history: make([][2]int, 0),
		}
	}

	return Rope{
		knots: knots,
	}
}

func (r *Rope) Move(move Movement) {
	for i := 0; i < move.Amount; i++ {
		switch move.Direction {
		case RightMovement:
			r.knots[0].moveRight()
		case LeftMovement:
			r.knots[0].moveLeft()
		case UpMovement:
			r.knots[0].moveUp()
		case DownMovement:
			r.knots[0].moveDown()
		default:
			panic("grid: unsupported movement direction")
		}

		for i := 1; i < knotCount; i++ {
			r.knots[i].follow(r.knots[i-1])
		}
	}
}

func (r *Rope) CountTailUniqueLocations() int {
	positions := map[[2]int]int{}
	last := r.knots[knotCount-1]

	positions[[2]int{last.x, last.y}] += 1

	for _, position := range last.history {
		positions[position] += 1
	}

	return len(positions)
}
