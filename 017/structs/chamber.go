package structs

import (
	"fmt"
	"strings"
)

const (
	chamberWidth             = 7
	chamberSpanPointX        = 2
	chamberSpanPointYOffset  = 3
	chamberRockStartCapacity = 5_000

	rockPatterCount = 5
)

type location struct {
	x, y int64
}

type Chamber struct {
	jets         jet
	rocks        map[location]bool
	rockCount    uint
	highestPoint int64
	spanPoint    int64
	width        int64
}

func NewChamber() *Chamber {
	return &Chamber{
		rocks:        make(map[location]bool, chamberRockStartCapacity),
		rockCount:    0,
		highestPoint: 0,
		spanPoint:    chamberSpanPointX,
		width:        chamberWidth,
	}
}

func (c *Chamber) SetJets(jets string) {
	c.jets = jet{}
	c.jets.setPattern([]rune(jets))
}

func (c *Chamber) HowManyUnitsTallWillTheTowerOfRocksBeAfterNRocksHaveStoppedFalling(limit uint) int64 {
	var (
		falling   *rock
		moveLeft  = '<'
		moveRight = '>'
	)

	for i := uint64(0); c.rockCount < limit; i++ {
		if falling == nil {
			falling = &rock{}
			spanAt := location{x: c.spanPoint, y: c.highestPoint + chamberSpanPointYOffset}
			falling.spanAt(c.rockCount%rockPatterCount, spanAt)
		}

		move := c.jets.getNextDirection(i)

		if move == moveLeft && falling.left > 0 && falling.canMoveLeft(c.rocks) {
			falling.moveLeft()
		}

		if move == moveRight && falling.right < chamberWidth-1 && falling.canMoveRight(c.rocks) {
			falling.moveRight()
		}

		if falling.bottom > 0 && falling.canFall(c.rocks) {
			falling.fall()
			continue
		}

		if falling.top >= c.highestPoint {
			c.highestPoint = falling.top + 1
		}

		c.rocks = falling.stop(c.rocks)

		c.rockCount++

		falling = nil
	}

	return c.highestPoint
}

func (c *Chamber) Print() {
	builder := strings.Builder{}

	for i := c.highestPoint; i >= 0; i-- {
		builder.WriteRune('|')
		for j := int64(0); j < chamberWidth; j++ {
			if c.rocks[location{x: j, y: i}] {
				builder.WriteRune('#')
				continue
			}
			builder.WriteRune('.')
		}
		builder.WriteRune('|')
		fmt.Println(builder.String())
		builder.Reset()
	}

	fmt.Println("+-------+")
}

func (c *Chamber) PrintFalling(falling *rock) {
	builder := strings.Builder{}

	start := c.highestPoint

	if falling.top > start {
		start = falling.top
	}

	for i := start; i >= 0; i-- {
		builder.WriteRune('|')
		for j := int64(0); j < chamberWidth; j++ {
			if c.rocks[location{x: j, y: i}] || falling.contains(j, i) {
				if falling.contains(j, i) {
					builder.WriteRune('@')
				} else {
					builder.WriteRune('#')
				}

				continue
			}
			builder.WriteRune('.')
		}
		builder.WriteRune('|')
		fmt.Println(builder.String())
		builder.Reset()
	}

	fmt.Println("+-------+")
}
