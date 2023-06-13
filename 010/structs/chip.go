package structs

import (
	"errors"
	"strings"
)

type Chip struct {
	x                  int
	instructions       queue
	currentInstruction instruction
	cycle              uint
}

type ProbeResult struct {
	Cycle uint
	X     int
}

func NewChip() Chip {
	return Chip{
		x: 1,
	}
}

func (c *Chip) ParseAndQueueInstruction(input string) error {
	inst := strings.Split(input, " ")

	switch inst[0] {
	case "addx":
		ins, err := newAddXInstruction(inst[1])
		if err != nil {
			return err
		}
		c.instructions.enqueue(&ins)
	case "noop":
		ins := newNoOpInstruction()
		c.instructions.enqueue(&ins)
	default:
		return errors.New("unsupported instruction")
	}

	return nil
}

func (c *Chip) RunProgram(probeFunc func(uint) bool) <-chan ProbeResult {
	out := make(chan ProbeResult, c.instructions.size())

	go func(out chan<- ProbeResult, probeFunc func(uint) bool) {
		ins, err := c.instructions.dequeue()
		if err != nil {
			panic(err)
		}

		c.currentInstruction = ins

		for c.currentInstruction != nil {
			c.currentInstruction.Tick()
			c.cycle++

			if probeFunc(c.cycle) {
				out <- ProbeResult{
					Cycle: c.cycle,
					X:     c.x,
				}
			}

			if c.currentInstruction.ChangeState(c) {
				c.currentInstruction = nil
			}

			if c.currentInstruction == nil && c.instructions.size() > 0 {
				c.currentInstruction, err = c.instructions.dequeue()
				if err != nil {
					panic(err)
				}
			}
		}

		close(out)
	}(out, probeFunc)

	return out
}

func (c *Chip) ClearInstructions() {
	c.instructions.clear()
}
