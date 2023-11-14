package encrypted

const (
	defaultFileSize = 5000
)

type node struct {
	moves int64
	left  *node
	right *node
}

type File struct {
	items []*node
	key   int64
}

func New(key int64) File {
	return File{
		items: make([]*node, 0, defaultFileSize),
		key:   key,
	}
}

func (f *File) AddItem(item int64) {
	n := &node{moves: item * f.key}

	if len(f.items) == 0 {
		n.right = n
		n.left = n
	} else {
		n.right = f.items[0]
		n.left = f.items[len(f.items)-1]

		f.items[0].left = n
		f.items[len(f.items)-1].right = n
	}

	f.items = append(f.items, n)
}

func (f *File) Size() int {
	return len(f.items)
}

func (f *File) GetCoordinates(coordinates ...int) int64 {
	var (
		result int64
		zero   *node

		m = int64(len(f.items) - 1)
	)

	for i := 0; i < len(f.items); i++ {
		var (
			current = f.items[i]
		)

		if current.moves == 0 {
			zero = current
			continue
		}

		p := current

		if current.moves > 0 {
			for i := int64(0); i < current.moves%m; i++ {
				p = p.right
			}
		} else {
			for i := int64(0); i < -(current.moves-1)%m; i++ {
				p = p.left
			}
		}

		if current == p {
			continue
		}

		current.right.left = current.left
		current.left.right = current.right

		p.right.left = current
		current.right = p.right

		p.right = current
		current.left = p
	}

	for i := 0; i < len(coordinates); i++ {
		current := zero

		for j := 0; j < coordinates[i]; j++ {
			current = current.right
		}

		result += current.moves
	}

	return result
}
