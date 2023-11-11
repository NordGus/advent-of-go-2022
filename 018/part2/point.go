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

func (p *point) cool() {
	if p.material == lava || p.material == air {
		return
	}

	if len(p.neighbors) == 0 {
		return
	}

	for _, neighbor := range p.neighbors {
		if neighbor.material == air && p.material != lava {
			p.material = air
			break
		}
	}

	for _, neighbor := range p.neighbors {
		neighbor.cool()
	}
}
