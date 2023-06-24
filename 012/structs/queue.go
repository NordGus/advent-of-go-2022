package structs

type queueNode struct {
	item     *node
	distance int
}

type queue struct {
	items []queueNode
}

func (pq *queue) enqueue(item *node, distance int) {
	pn := queueNode{item: item, distance: distance}

	pq.items = append(pq.items, pn)
}

func (pq *queue) dequeue() (*node, int) {
	if len(pq.items) == 0 {
		return nil, 0
	}

	item := pq.items[0].item
	distance := pq.items[0].distance

	pq.items[0].item = nil

	if len(pq.items) == 1 {
		pq.items = make([]queueNode, 0)
	} else {
		pq.items = pq.items[1:len(pq.items)]
	}

	return item, distance
}

func (pq *queue) size() int {
	return len(pq.items)
}

func (pq *queue) clear() {
	for {
		value, _ := pq.dequeue()
		if value == nil {
			break
		}
	}
}
