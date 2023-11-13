package factory

const (
	initialRobotCapacity = 5
)

type Blueprint struct {
	id int

	robots map[Resource]robot
}

func NewBlueprint(id int) Blueprint {
	return Blueprint{
		id:     id,
		robots: make(map[Resource]robot, initialRobotCapacity),
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

	b.robots[rbt.Resource] = rbt

	return err
}
