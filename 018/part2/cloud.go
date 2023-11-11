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
		for _, neighbor := range c.cloud[lv.x][lv.y][lv.z].neighbors {
			if neighbor.material == lava || neighbor.material == steam {
				sides -= 1
			}
		}
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
				cloud[i][j][k] = newPoint(i, j, k, steam)
			}
		}
	}

	for _, lv := range c.lava {
		cloud[lv.x][lv.y][lv.z].material = lv.material
	}

	c.cloud = cloud

	for i := 0; i < c.maxX; i++ {
		for j := 0; j < c.maxY; j++ {
			for k := 0; k < c.maxZ; k++ {
				c.cloud[i][j][k].neighbors = append(
					c.cloud[i][j][k].neighbors,
					c.getNeighbor(i-1, j, k),
					c.getNeighbor(i+1, j, k),
					c.getNeighbor(i, j-1, k),
					c.getNeighbor(i, j+1, k),
					c.getNeighbor(i, j, k-1),
					c.getNeighbor(i, j, k+1),
				)
			}
		}
	}

	c.cloud[0][0][0].cool()

	c.built = true
}

func (c *Cloud) getNeighbor(x int, y int, z int) *point {
	if x < 0 || y < 0 || z < 0 || x >= c.maxX || y >= c.maxY || z >= c.maxZ {
		pnt := newPoint(x, y, z, air)

		return &pnt
	}

	return &c.cloud[x][y][z]
}
