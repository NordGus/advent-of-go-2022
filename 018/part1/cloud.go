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

func (c *Cloud) CountSidesThatAreNotConnectedBetweenCubes() int {
	var (
		space = c.buildSpace()
		sides = len(c.points) * shapeSidesCount
	)

	for _, pnt := range c.points {
		neighbors := make([]point, 0, shapeSidesCount)

		negXNeighbor := getNeighbor(space, pnt.planeX()-1, pnt.planeY(), pnt.planeZ())
		posXNeighbor := getNeighbor(space, pnt.planeX()+1, pnt.planeY(), pnt.planeZ())

		negYNeighbor := getNeighbor(space, pnt.planeX(), pnt.planeY()-1, pnt.planeZ())
		posYNeighbor := getNeighbor(space, pnt.planeX(), pnt.planeY()+1, pnt.planeZ())

		negZNeighbor := getNeighbor(space, pnt.planeX(), pnt.planeY(), pnt.planeZ()-1)
		posZNeighbor := getNeighbor(space, pnt.planeX(), pnt.planeY(), pnt.planeZ()+1)

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
		space[p.planeX()][p.planeY()][p.planeZ()] = p
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
