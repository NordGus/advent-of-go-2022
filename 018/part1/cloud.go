package part1

const (
	startingPointsSize = 50
	shapeSidesCount    = 6
)

type Cloud struct {
	points []point
	maxX   int
	maxY   int
	maxZ   int
}

func NewCloud() *Cloud {
	return &Cloud{
		points: make([]point, 0, startingPointsSize),
	}
}

func (c *Cloud) AddPoint(data [3]int) {
	p := newPoint(data[0], data[1], data[2])

	if p.x >= c.maxX {
		c.maxX = p.x + 1
	}

	if p.y >= c.maxY {
		c.maxY = p.y + 1
	}

	if p.z >= c.maxZ {
		c.maxZ = p.z + 1
	}

	c.points = append(c.points, p)
}

func (c *Cloud) CountSidesThatAreNotConnectedBetweenCubes() int {
	var (
		space = c.buildSpace()
		sides = len(c.points) * shapeSidesCount
	)

	for _, pnt := range c.points {
		neighbors := make([]point, 0, shapeSidesCount)

		negXNeighbor := getNeighbor(space, pnt.x-1, pnt.y, pnt.z)
		posXNeighbor := getNeighbor(space, pnt.x+1, pnt.y, pnt.z)

		negYNeighbor := getNeighbor(space, pnt.x, pnt.y-1, pnt.z)
		posYNeighbor := getNeighbor(space, pnt.x, pnt.y+1, pnt.z)

		negZNeighbor := getNeighbor(space, pnt.x, pnt.y, pnt.z-1)
		posZNeighbor := getNeighbor(space, pnt.x, pnt.y, pnt.z+1)

		if negXNeighbor.active {
			neighbors = append(neighbors, negXNeighbor)
		}

		if posXNeighbor.active {
			neighbors = append(neighbors, posXNeighbor)
		}

		if negYNeighbor.active {
			neighbors = append(neighbors, negYNeighbor)
		}

		if posYNeighbor.active {
			neighbors = append(neighbors, posYNeighbor)
		}

		if negZNeighbor.active {
			neighbors = append(neighbors, negZNeighbor)
		}

		if posZNeighbor.active {
			neighbors = append(neighbors, posZNeighbor)
		}

		sides -= len(neighbors)
	}

	return sides
}

func (c *Cloud) buildSpace() [][][]point {
	space := make([][][]point, c.maxX)

	for i := 0; i < c.maxX; i++ {
		space[i] = make([][]point, c.maxY)

		for j := 0; j < c.maxY; j++ {
			space[i][j] = make([]point, c.maxZ)
		}
	}

	for _, p := range c.points {
		space[p.x][p.y][p.z] = p
	}

	return space
}

func getNeighbor(space [][][]point, x int, y int, z int) point {
	if x < 0 || y < 0 || z < 0 {
		return point{}
	}

	if x >= len(space) || y >= len(space[0]) || z >= len(space[0][0]) {
		return point{}
	}

	return space[x][y][z]
}
