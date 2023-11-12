package part1

import (
	"errors"
	"fmt"
)

const (
	initialRobotCostCapacity = 3
)

type robot struct {
	Resource Resource

	Cost map[Resource]int
}

func newRobot(resource string, costs map[string]int) (robot, error) {
	var (
		rsrc Resource
		cst  = make(map[Resource]int, initialRobotCostCapacity)

		rbt robot
		err error
	)

	rsrc, err = getResource(resource)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("blueprint: %s is not a valid resource", resource))

		return rbt, err
	}

	for res, cost := range costs {
		var r Resource

		r, err = getResource(res)
		if err != nil {
			err = errors.Join(err, fmt.Errorf("blueprint: %s is not a valid resource", res))

			return rbt, err
		}

		cst[r] = cost
	}

	rbt.Resource = rsrc
	rbt.Cost = cst

	return rbt, err
}
