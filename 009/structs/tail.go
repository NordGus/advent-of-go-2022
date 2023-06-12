package structs

type tail struct {
	x, y    int
	history [][2]int
}

func (t *tail) savePosition() {
	t.history = append(t.history, [2]int{t.x, t.y})
}

func (t *tail) move(h *head) {
	t.savePosition()

	if h.x >= t.x+2 && h.y > t.y {
		t.x++
		t.y++
	}

	if h.x >= t.x+2 && h.y < t.y {
		t.x++
		t.y--
	}

	if h.x <= t.x-2 && h.y > t.y {
		t.x--
		t.y++
	}

	if h.x <= t.x-2 && h.y < t.y {
		t.x--
		t.y--
	}

	if h.y >= t.y+2 && h.x > t.x {
		t.x++
		t.y++
	}

	if h.y >= t.y+2 && h.x < t.x {
		t.x--
		t.y++
	}

	if h.y <= t.y-2 && h.x > t.x {
		t.x++
		t.y--
	}

	if h.y <= t.y-2 && h.x < t.x {
		t.x--
		t.y--
	}

	if h.x > t.x+1 {
		t.x++
	}

	if h.y > t.y+1 {
		t.y++
	}

	if h.x < t.x-1 {
		t.x--
	}

	if h.y < t.y-1 {
		t.y--
	}
}
