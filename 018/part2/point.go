package part2

type materialType int

const (
	air materialType = iota
	steam
	lava
)

type point struct {
	x, y, z   int
	material  materialType
	neighbors []*point
}

func newPoint(x, y, z int, material materialType) point {
	return point{
		x:         x,
		y:         y,
		z:         z,
		material:  material,
		neighbors: make([]*point, 0, shapeSidesCount),
	}
}
