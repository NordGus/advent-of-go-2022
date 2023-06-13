package structs

import (
	"errors"
	"strings"
)

const (
	ScreenHorizontalResolution = 40
	ScreenVerticalResolution   = 6
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

func (c *Chip) RenderScreen() <-chan string {
	out := make(chan string, ScreenVerticalResolution)

	go func(out chan<- string) {
		const brightPixel = "#"
		const darkPixel = "."

		scanLine := strings.Builder{}
		s := sprite{pos: c.x}

		ins, err := c.instructions.dequeue()
		if err != nil {
			panic(err)
		}

		c.currentInstruction = ins
		scanLine.WriteString(brightPixel)

		for c.currentInstruction != nil {
			c.currentInstruction.Tick()
			c.cycle++

			if c.currentInstruction.ChangeState(c) {
				c.currentInstruction = nil
			}

			s.move(c.x)

			if s.isVisible(int(c.cycle % ScreenHorizontalResolution)) {
				scanLine.WriteString(brightPixel)
			} else {
				scanLine.WriteString(darkPixel)
			}

			if c.currentInstruction == nil && c.instructions.size() > 0 {
				c.currentInstruction, err = c.instructions.dequeue()
				if err != nil {
					panic(err)
				}
			}

			if len(scanLine.String()) == ScreenHorizontalResolution {
				out <- scanLine.String()
				scanLine.Reset()
			}
		}

		close(out)
	}(out)

	return out
}

func (c *Chip) ClearInstructions() {
	c.instructions.clear()
}
