package structs

type line struct {
	points      []point
	lowestPoint int
}

func (l *line) detectCollision(p *point) bool {
	buffer := 0.1

	if p.y > l.lowestPoint {
		return false
	}

	for i := 0; i+1 < len(l.points); i++ {
		d1 := p.distance(&l.points[i])
		d2 := p.distance(&l.points[i+1])
		lineLen := l.points[i].distance(&l.points[i+1])

		if d1+d2 >= lineLen-buffer && d1+d2 <= lineLen+buffer {
			return true
		}
	}

	return false
}
