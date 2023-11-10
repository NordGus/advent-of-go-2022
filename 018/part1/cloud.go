package part1

const (
	startingPointsSize = 50
)

type Cloud struct {
	points []point
	maxX   int
	maxY   int
	maxZ   int
}

type point struct {
	x, y, z int
}

func NewCloud() *Cloud {
	return &Cloud{
		points: make([]point, 0, startingPointsSize),
	}
}

func (c *Cloud) AddPoint(data [3]int) {
	p := point{x: data[0], y: data[1], z: data[2]}

	if p.x > c.maxX {
		c.maxX = p.x
	}

	if p.y > c.maxY {
		c.maxY = p.y
	}

	if p.z > c.maxZ {
		c.maxZ = p.z
	}

	c.points = append(c.points, p)
}
