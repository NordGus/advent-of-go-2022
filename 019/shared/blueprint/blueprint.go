package blueprint

const (
	initialRobotCapacity = 5
)

type Blueprint struct {
	ID int

	Robots []robot
}

func New(id int) Blueprint {
	return Blueprint{
		ID:     id,
		Robots: make([]robot, 0, initialRobotCapacity),
	}
}

func (b *Blueprint) AddRobotRecipe(resource string, costs map[string]int) error {
	var (
		rbt robot
		err error
	)

	rbt, err = newRobot(resource, costs)
	if err != nil {
		return err
	}

	b.Robots = append(b.Robots, rbt)

	return err
}
