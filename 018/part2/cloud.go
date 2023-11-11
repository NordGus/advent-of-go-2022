package part2

const (
	startingPointsSize = 50
	shapeSidesCount    = 6
)

type Cloud struct {
	lava  []point
	cloud [][][]point
	maxX  int
	maxY  int
	maxZ  int

	built bool
}

func NewCloud() *Cloud {
	return &Cloud{
		lava: make([]point, 0, startingPointsSize),
	}
}

func (c *Cloud) AddLavaPoint(data [3]int) {
	p := newPoint(data[0], data[1], data[2], lava)

	if p.x >= c.maxX {
		c.maxX = p.x + 1
	}

	if p.y >= c.maxY {
		c.maxY = p.y + 1
	}

	if p.z >= c.maxZ {
		c.maxZ = p.z + 1
	}

	c.lava = append(c.lava, p)
}

func (c *Cloud) CountExternalSurfaceAreaOfLavaDroplet() int {
	var (
		sides = len(c.lava) * shapeSidesCount
	)

	c.buildCloud()

	for _, lv := range c.lava {
		neighbors := make([]point, 0, shapeSidesCount)

		negXNeighbor := c.getNeighbor(lv.x-1, lv.y, lv.z)
		posXNeighbor := c.getNeighbor(lv.x+1, lv.y, lv.z)

		negYNeighbor := c.getNeighbor(lv.x, lv.y-1, lv.z)
		posYNeighbor := c.getNeighbor(lv.x, lv.y+1, lv.z)

		negZNeighbor := c.getNeighbor(lv.x, lv.y, lv.z-1)
		posZNeighbor := c.getNeighbor(lv.x, lv.y, lv.z+1)

		if negXNeighbor.material == lava || negXNeighbor.material == steam {
			neighbors = append(neighbors, negXNeighbor)
		}

		if posXNeighbor.material == lava || posXNeighbor.material == steam {
			neighbors = append(neighbors, posXNeighbor)
		}

		if negYNeighbor.material == lava || negYNeighbor.material == steam {
			neighbors = append(neighbors, negYNeighbor)
		}

		if posYNeighbor.material == lava || posYNeighbor.material == steam {
			neighbors = append(neighbors, posYNeighbor)
		}

		if negZNeighbor.material == lava || negZNeighbor.material == steam {
			neighbors = append(neighbors, negZNeighbor)
		}

		if posZNeighbor.material == lava || posZNeighbor.material == steam {
			neighbors = append(neighbors, posZNeighbor)
		}

		sides -= len(neighbors)
	}

	return sides
}

func (c *Cloud) buildCloud() {
	if c.built {
		return
	}

	cloud := make([][][]point, c.maxX)

	for i := 0; i < c.maxX; i++ {
		cloud[i] = make([][]point, c.maxY)

		for j := 0; j < c.maxY; j++ {
			cloud[i][j] = make([]point, c.maxZ)

			for k := 0; k < c.maxZ; k++ {
				cloud[i][j][k] = newPoint(i, j, k, air)
			}
		}
	}

	for _, p := range c.lava {
		cloud[p.x][p.y][p.z].material = p.material
	}

	c.cloud = cloud
	c.built = true
}

func (c *Cloud) getNeighbor(x int, y int, z int) point {
	if x < 0 || y < 0 || z < 0 || x >= c.maxX || y >= c.maxY || z >= c.maxZ {
		return newPoint(x, y, z, air)
	}

	return c.cloud[x][y][z]
}
