package structs

type Cave struct {
	sandSource  point
	rocks       []line
	sand        map[point]bool
	lowestPoint int
}

func NewCave() Cave {
	return Cave{
		sandSource: point{x: 500, y: 0},
		rocks:      make([]line, 0, 1_000),
		sand:       make(map[point]bool, 1_000),
	}
}

func (c *Cave) AddLine(points ...[2]int) {
	l := line{
		points: make([]point, len(points)),
	}

	for i := 0; i < len(points); i++ {
		l.points[i] = point{
			x: points[i][0],
			y: points[i][1],
		}

		if l.points[i].y > c.lowestPoint {
			c.lowestPoint = l.points[i].y
		}

		if l.points[i].y > l.lowestPoint {
			l.lowestPoint = l.points[i].y
		}
	}

	c.rocks = append(c.rocks, l)
}

func (c *Cave) HowManyUnitsOfSandBeforeOverflowing() int {
	grain := c.sandSource

	for grain.y <= c.lowestPoint {
		if c.sandCanFallDown(grain) {
			grain = point{x: grain.x, y: grain.y + 1}
			continue
		}

		if c.sandCanFallToTheLeft(grain) {
			grain = point{x: grain.x - 1, y: grain.y + 1}
			continue
		}

		if c.sandCanFallToTheRight(grain) {
			grain = point{x: grain.x + 1, y: grain.y + 1}
			continue
		}

		c.sand[grain] = true
		grain = c.sandSource
	}

	return len(c.sand)
}

func (c *Cave) sandCanFallDown(grain point) bool {
	sensor := point{x: grain.x, y: grain.y + 1}

	if c.sand[sensor] {
		return false
	}

	for i := 0; i < len(c.rocks); i++ {
		if c.rocks[i].detectCollision(&sensor) {
			return false
		}
	}

	return true
}

func (c *Cave) sandCanFallToTheLeft(grain point) bool {
	sensor := point{x: grain.x - 1, y: grain.y + 1}

	if c.sand[sensor] {
		return false
	}

	for i := 0; i < len(c.rocks); i++ {
		if c.rocks[i].detectCollision(&sensor) {
			return false
		}
	}

	return true
}

func (c *Cave) sandCanFallToTheRight(grain point) bool {
	sensor := point{x: grain.x + 1, y: grain.y + 1}

	if c.sand[sensor] {
		return false
	}

	for i := 0; i < len(c.rocks); i++ {
		if c.rocks[i].detectCollision(&sensor) {
			return false
		}
	}

	return true
}
