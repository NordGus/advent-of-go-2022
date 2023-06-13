package structs

type instruction interface {
	ChangeState(c *Chip) bool
	Tick()
}
