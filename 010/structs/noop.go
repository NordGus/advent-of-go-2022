package structs

type noop struct {
	ticks int
}

func newNoOpInstruction() noop {
	return noop{ticks: 1}
}

func (i *noop) ChangeState(c *Chip) bool {
	return i.ticks <= 0
}

func (i *noop) Tick() {
	i.ticks--
}
