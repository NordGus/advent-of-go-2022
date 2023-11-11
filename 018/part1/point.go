package part1

type point struct {
	x, y, z int
	active  bool
}

func newPoint(x, y, z int) point {
	return point{x: x, y: y, z: z, active: true}
}
