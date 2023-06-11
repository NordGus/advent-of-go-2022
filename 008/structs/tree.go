package structs

type Tree struct {
	xCoordinate int
	yCoordinate int
	height      int
	scenicScore int
}

func (t *Tree) IsTaller(neighbor Tree) bool {
	return t.height > neighbor.height
}
