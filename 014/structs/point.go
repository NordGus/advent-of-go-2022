package structs

import "math"

type point struct {
	x, y int
}

func (p *point) distance(p2 *point) float64 {
	return math.Sqrt((math.Pow(float64(p2.x-p.x), 2) + math.Pow(float64(p2.y-p.y), 2)))
}
