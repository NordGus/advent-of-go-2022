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
	rockCacheCapacity        = 100

	rockPatterCount int64 = 5

	moveLeft  = '<'
	moveRight = '>'

	simSaveStateThreshold int64 = 10
)

type location struct {
	x, y   int64
	isRock bool
}

type state struct {
	heighDiff         int64
	innerHeight       int64
	rockPatternIndex  int64
	jetDirectionIndex uint64
	positionMask1     uint64
	positionMask2     uint64
	positionMask3     uint64
	positionMask4     uint64
}

type stateData struct {
	rockCount    int64
	highestPoint int64
}

type Chamber struct {
	jets         jet
	collisions   map[location]bool
	states       map[state]stateData
	lastState    state
	rocks        []rock
	rockCount    int64
	highestPoint int64
	spanPoint    int64
	width        int64
}

func NewChamber() *Chamber {
	return &Chamber{
		collisions:   make(map[location]bool, chamberRockStartCapacity),
		rocks:        make([]rock, 0, rockCacheCapacity),
		states:       make(map[state]stateData, rockCacheCapacity),
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

func (c *Chamber) HowManyUnitsTallWillTheTowerOfRocksBeAfterNRocksHaveStoppedFalling(limit int64) int64 {
	var falling *rock

	for i := uint64(0); c.rockCount < limit; i++ {
		if falling == nil && c.rockCount >= simSaveStateThreshold {
			if found, cached, current := c.findPatternOrSaveSimState(i); found {
				lastCachedRock := c.rocks[cached.rockCount-1]
				lastCurrentRock := c.rocks[current.rockCount-1]

				cachedRocks := c.rocks[cached.rockCount-simSaveStateThreshold : cached.rockCount]
				currentRocks := c.rocks[current.rockCount-simSaveStateThreshold : current.rockCount]

				fmt.Printf("cached \n %+v \n %+v \n", lastCachedRock, cachedRocks[len(cachedRocks)-1])
				fmt.Printf("current \n %+v \n %+v \n", lastCurrentRock, currentRocks[len(currentRocks)-1])

				cachedCollisions := make(map[location]bool, chamberRockStartCapacity)
				currentCollisions := make(map[location]bool, chamberRockStartCapacity)

				for i := int64(0); i < simSaveStateThreshold; i++ {
					cachedCollisions = cachedRocks[i].stop(cachedCollisions)
					currentCollisions = currentRocks[i].stop(currentCollisions)
				}

				builder := strings.Builder{}
				count := 0

				for i := lastCachedRock.top; i >= cachedRocks[0].bottom; i-- {
					for j := int64(0); j < c.width; j++ {
						if cachedCollisions[location{x: j, y: i, isRock: true}] {
							builder.WriteRune('#')
							continue
						}
						builder.WriteRune('.')
					}

					fmt.Println(builder.String())
					builder.Reset()
					count++
				}

				fmt.Println("lines:", count)
				count = 0

				for i := lastCurrentRock.top; i >= currentRocks[0].bottom; i-- {
					for j := int64(0); j < c.width; j++ {
						if currentCollisions[location{x: j, y: i, isRock: true}] {
							builder.WriteRune('#')
							continue
						}
						builder.WriteRune('.')
					}

					fmt.Println(builder.String())
					builder.Reset()
					count++
				}

				fmt.Println("lines:", count)

				fmt.Println(lastCachedRock.top - cachedRocks[0].bottom + 1)
				fmt.Println(lastCurrentRock.top - currentRocks[0].bottom + 1)

				return c.highestPoint
			}
		}

		if falling == nil {
			falling = &rock{}
			spanAt := location{x: c.spanPoint, y: c.highestPoint + chamberSpanPointYOffset}
			falling.spanAt(c.rockCount%rockPatterCount, spanAt)
		}

		move := c.jets.getNextDirection(i)

		if move == moveLeft && falling.left > 0 && falling.canMoveLeft(c.collisions) {
			falling.moveLeft()
		}

		if move == moveRight && falling.right < c.width-1 && falling.canMoveRight(c.collisions) {
			falling.moveRight()
		}

		if falling.bottom > 0 && falling.canFall(c.collisions) {
			falling.fall()
			continue
		}

		if falling.top >= c.highestPoint {
			c.highestPoint = falling.top + 1
		}

		c.collisions = falling.stop(c.collisions)
		c.rocks = append(c.rocks, *falling)

		c.rockCount++

		falling = nil
	}

	return c.highestPoint
}

func (c *Chamber) findPatternOrSaveSimState(tick uint64) (bool, stateData, stateData) {
	var (
		s    state
		data stateData

		rocks = c.rocks[c.rockCount-simSaveStateThreshold:]
		ls    = c.lastState
		cs    = make(map[location]bool, chamberRockStartCapacity)
	)

	if len(c.states) == 0 {
		s.heighDiff = c.highestPoint
	} else {
		s.heighDiff = c.highestPoint - c.states[ls].highestPoint
	}

	s.jetDirectionIndex = c.jets.getDirectionIndex(tick)
	s.rockPatternIndex = c.rockCount % rockPatterCount
	s.innerHeight = rocks[len(rocks)-1].top - rocks[0].bottom + 1

	data.highestPoint = c.highestPoint
	data.rockCount = c.rockCount

	for i := 0; i < len(rocks); i++ {
		cs = rocks[i].stop(cs)
	}

	top := rocks[len(rocks)-1].top
	bottom := rocks[0].bottom
	size := uint(0)

	for i := bottom; i <= top; i++ {
		for j := int64(0); j < c.width; j++ {
			if size >= 256 {
				continue
			}

			pos := cs[location{x: j, y: i, isRock: true}]

			if !pos {
				size++
				continue
			}

			mask := uint64(1)
			mask = mask << (size % 64)

			if size <= 64 {
				s.positionMask1 = s.positionMask1 | mask
			} else if size <= 128 {
				s.positionMask2 = s.positionMask2 | mask
			} else if size <= 192 {
				s.positionMask3 = s.positionMask3 | mask
			} else {
				s.positionMask4 = s.positionMask4 | mask
			}

			size++
		}
	}

	c.lastState = s

	oldData, ok := c.states[s]
	if ok {
		fmt.Printf("%+v: %+v \n", s, data)
		fmt.Printf("%+v: %+v \n", s, oldData)
		fmt.Printf("%+v: %+v \n", ls, c.states[ls])

		return true, oldData, data
	}

	c.states[s] = data

	return false, c.states[ls], data
}

func (c *Chamber) Print() {
	builder := strings.Builder{}

	for i := c.highestPoint; i >= 0; i-- {
		builder.WriteRune('|')
		for j := int64(0); j < c.width; j++ {
			if c.collisions[location{x: j, y: i}] {
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
		for j := int64(0); j < c.width; j++ {
			if c.collisions[location{x: j, y: i}] || falling.contains(j, i) {
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
