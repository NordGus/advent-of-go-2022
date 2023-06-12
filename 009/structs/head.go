package structs

type head struct {
	x, y    int
	history [][2]int
}

func (h *head) savePosition() {
	h.history = append(h.history, [2]int{h.x, h.y})
}

func (h *head) moveRight() {
	h.savePosition()
	h.x++
}

func (h *head) moveLeft() {
	h.savePosition()
	h.x--
}

func (h *head) moveUp() {
	h.savePosition()
	h.y++
}

func (h *head) moveDown() {
	h.savePosition()
	h.y--
}
