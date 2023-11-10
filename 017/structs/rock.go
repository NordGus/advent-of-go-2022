package structs

import "fmt"

type rock struct {
	rocks  [][]location
	left   int64
	top    int64
	right  int64
	bottom int64
}

func (r *rock) spanAt(index int64, span location) {
	r.setPattern(index)
	r.left = span.x
	r.right = int64(len(r.rocks[0])) + span.x - 1

	r.top = int64(len(r.rocks)) + span.y - 1
	r.bottom = span.y

	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			r.rocks[i][j].x += span.x
			r.rocks[i][j].y += span.y
		}
	}
}

func (r *rock) changeHeight(h int64) {
	oldTop := r.top
	r.top = h
	r.bottom = r.top - (oldTop - r.bottom)

	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			r.rocks[i][j].y = h - (oldTop - r.rocks[i][j].y)
		}
	}
}

func (r *rock) canMoveLeft(obstacles map[location]bool) bool {
	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			if !r.rocks[i][j].isRock {
				continue
			}

			if obstacles[location{x: r.rocks[i][j].x - 1, y: r.rocks[i][j].y, isRock: true}] {
				return false
			}
		}
	}

	return true
}

func (r *rock) canMoveRight(obstacles map[location]bool) bool {
	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			if !r.rocks[i][j].isRock {
				continue
			}

			if obstacles[location{x: r.rocks[i][j].x + 1, y: r.rocks[i][j].y, isRock: true}] {
				return false
			}
		}
	}

	return true
}

func (r *rock) canFall(obstacles map[location]bool) bool {
	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			if !r.rocks[i][j].isRock {
				continue
			}

			if obstacles[location{x: r.rocks[i][j].x, y: r.rocks[i][j].y - 1, isRock: true}] {
				return false
			}
		}
	}

	return true
}

func (r *rock) moveLeft() {
	r.left--
	r.right--

	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			r.rocks[i][j].x--
		}
	}
}

func (r *rock) moveRight() {
	r.left++
	r.right++

	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			r.rocks[i][j].x++
		}
	}
}

func (r *rock) fall() {
	r.top--
	r.bottom--

	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			r.rocks[i][j].y--
		}
	}
}

func (r *rock) stop(obstacles map[location]bool) map[location]bool {
	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			if !r.rocks[i][j].isRock {
				continue
			}

			obstacles[r.rocks[i][j]] = true
		}
	}

	return obstacles
}

func (r *rock) contains(x int64, y int64) bool {
	for i := 0; i < len(r.rocks); i++ {
		for j := 0; j < len(r.rocks[i]); j++ {
			if !r.rocks[i][j].isRock {
				continue
			}

			if r.rocks[i][j].x == x && r.rocks[i][j].y == y {
				return true
			}
		}
	}

	return false
}

func (r *rock) setPattern(index int64) {
	switch index {
	case 0:
		// ####
		r.rocks = [][]location{
			{location{x: 0, y: 0, isRock: true}, location{x: 1, y: 0, isRock: true}, location{x: 2, y: 0, isRock: true}, location{x: 3, y: 0, isRock: true}},
		}
	case 1:
		// .#.
		// ###
		// .#.
		r.rocks = [][]location{
			{location{x: 0, y: 0}, location{x: 1, y: 0, isRock: true}, location{x: 2, y: 0}},
			{location{x: 0, y: 1, isRock: true}, location{x: 1, y: 1, isRock: true}, location{x: 2, y: 1, isRock: true}},
			{location{x: 0, y: 2}, location{x: 1, y: 2, isRock: true}, location{x: 2, y: 2}},
		}
	case 2:
		// ..#
		// ..#
		// ###
		r.rocks = [][]location{
			{location{x: 0, y: 0, isRock: true}, location{x: 1, y: 0, isRock: true}, location{x: 2, y: 0, isRock: true}},
			{location{x: 0, y: 1}, location{x: 1, y: 1}, location{x: 2, y: 1, isRock: true}},
			{location{x: 0, y: 2}, location{x: 1, y: 2}, location{x: 2, y: 2, isRock: true}},
		}
	case 3:
		// #
		// #
		// #
		// #
		r.rocks = [][]location{
			{location{x: 0, y: 0, isRock: true}},
			{location{x: 0, y: 1, isRock: true}},
			{location{x: 0, y: 2, isRock: true}},
			{location{x: 0, y: 3, isRock: true}},
		}
	case 4:
		// ##
		// ##
		r.rocks = [][]location{
			{location{x: 0, y: 0, isRock: true}, location{x: 1, y: 0, isRock: true}},
			{location{x: 0, y: 1, isRock: true}, location{x: 1, y: 1, isRock: true}},
		}
	default:
		panic(fmt.Sprintf("there's something wrong with the math: %v", index))
	}
}
