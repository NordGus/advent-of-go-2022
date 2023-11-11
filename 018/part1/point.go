package part1

type point struct {
	x, y, z int
	active  bool
}

func newPoint(x, y, z int) point {
	return point{x: x, y: y, z: z, active: true}
}

func (p *point) planeX() int {
	return p.x - 1
}

func (p *point) planeY() int {
	return p.y - 1
}

func (p *point) planeZ() int {
	return p.z - 1
}
