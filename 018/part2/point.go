package part2

type materialType int

const (
	air materialType = iota
	steam
	lava
)

type point struct {
	x, y, z  int
	material materialType
}

func newPoint(x, y, z int) point {
	return point{x: x, y: y, z: z}
}
