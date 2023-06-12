package structs

type knot struct {
	x, y    int
	history [][2]int
}

func (k *knot) savePosition() {
	k.history = append(k.history, [2]int{k.x, k.y})
}

func (k *knot) moveRight() {
	k.savePosition()
	k.x++
}

func (k *knot) moveLeft() {
	k.savePosition()
	k.x--
}

func (k *knot) moveUp() {
	k.savePosition()
	k.y++
}

func (k *knot) moveDown() {
	k.savePosition()
	k.y--
}

func (follower *knot) follow(leader *knot) {
	follower.savePosition()

	if leader.x >= follower.x+2 && leader.y > follower.y {
		follower.x++
		follower.y++
	}

	if leader.x >= follower.x+2 && leader.y < follower.y {
		follower.x++
		follower.y--
	}

	if leader.x <= follower.x-2 && leader.y > follower.y {
		follower.x--
		follower.y++
	}

	if leader.x <= follower.x-2 && leader.y < follower.y {
		follower.x--
		follower.y--
	}

	if leader.y >= follower.y+2 && leader.x > follower.x {
		follower.x++
		follower.y++
	}

	if leader.y >= follower.y+2 && leader.x < follower.x {
		follower.x--
		follower.y++
	}

	if leader.y <= follower.y-2 && leader.x > follower.x {
		follower.x++
		follower.y--
	}

	if leader.y <= follower.y-2 && leader.x < follower.x {
		follower.x--
		follower.y--
	}

	if leader.x > follower.x+1 {
		follower.x++
	}

	if leader.y > follower.y+1 {
		follower.y++
	}

	if leader.x < follower.x-1 {
		follower.x--
	}

	if leader.y < follower.y-1 {
		follower.y--
	}
}
