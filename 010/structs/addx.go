package structs

import (
	"strconv"
)

type addx struct {
	v     int
	ticks int
}

func newAddXInstruction(value string) (addx, error) {
	val, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return addx{}, err
	}

	return addx{v: int(val), ticks: 2}, nil
}

func (i *addx) ChangeState(c *Chip) bool {
	if i.ticks > 0 {
		return false
	}

	c.x += i.v

	return true
}

func (i *addx) Tick() {
	i.ticks--
}
