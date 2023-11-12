package part1

const (
	initialRobotCapacity = 5
)

type Blueprint struct {
	id int

	robots []robot
}

func NewBlueprint(id int) Blueprint {
	return Blueprint{
		id:     id,
		robots: make([]robot, 0, initialRobotCapacity),
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

	b.robots = append(b.robots, rbt)

	return err
}