package structs

import "sync"

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
		sand:       make(map[point]bool, 10_000),
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

		if i == 0 {
			l.highestPoint = l.points[i].y
			l.lowestPoint = l.points[i].y
			l.westernPoint = l.points[i].x
			l.easternPoint = l.points[i].x
		}

		if l.points[i].y > c.lowestPoint {
			c.lowestPoint = l.points[i].y
		}

		if l.points[i].y > l.lowestPoint {
			l.lowestPoint = l.points[i].y
		}

		if l.points[i].y < l.highestPoint {
			l.highestPoint = l.points[i].y
		}

		if l.points[i].x > l.easternPoint {
			l.easternPoint = l.points[i].x
		}

		if l.points[i].x < l.westernPoint {
			l.westernPoint = l.points[i].x
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

func (c *Cave) HowManyUnitsOfSandBeforeBlockage() int {
	grain := c.sandSource

	for {
		if grain.y == c.lowestPoint+1 {
			c.sand[grain] = true
			grain = c.sandSource
			continue
		}

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

		if grain == c.sandSource {
			c.sand[grain] = true
			break
		}

		c.sand[grain] = true
		grain = c.sandSource
	}

	return len(c.sand)
}

func (c *Cave) sandCanFallDown(grain point) bool {
	sensor := point{x: grain.x, y: grain.y + 1}

	return !c.pointCollides(sensor)
}

func (c *Cave) sandCanFallToTheLeft(grain point) bool {
	sensor := point{x: grain.x - 1, y: grain.y + 1}

	return !c.pointCollides(sensor)
}

func (c *Cave) sandCanFallToTheRight(grain point) bool {
	sensor := point{x: grain.x + 1, y: grain.y + 1}

	return !c.pointCollides(sensor)
}

func (c *Cave) pointCollides(p point) bool {
	if c.sand[p] {
		return true
	}

	collided := false
	collisions := make(chan bool, len(c.rocks)/2)
	var wg sync.WaitGroup

	for i := 0; i < len(c.rocks); i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup, out chan<- bool, l *line, p *point) {
			defer wg.Done()
			out <- l.detectCollision(p)
		}(&wg, collisions, &c.rocks[i], &p)
	}

	go func(wg *sync.WaitGroup, out chan<- bool) {
		wg.Wait()
		close(out)
	}(&wg, collisions)

	for collision := range collisions {
		if collision {
			collided = true
		}
	}

	return collided
}
