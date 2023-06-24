package structs

type node struct {
	x      int
	y      int
	height int
	edges  []*node
}

func (n *node) addEdge(to *node) {
	if to == nil {
		return
	}

	if n.height < to.height-1 {
		return
	}

	n.edges = append(n.edges, to)
}
