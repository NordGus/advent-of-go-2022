package encrypted

const (
	defaultFileSize = 5000
)

type node struct {
	moves int64
	next  *node
	prev  *node
}

type File struct {
	items  []*node
	zeroAt int
}

func New() File {
	return File{
		items:  make([]*node, 0, defaultFileSize),
		zeroAt: -1,
	}
}

func (f *File) AddItem(item int64) {
	n := &node{moves: item}

	if len(f.items) == 0 {
		n.next = n
		n.prev = n
	} else {
		n.next = f.items[0]
		n.prev = f.items[len(f.items)-1]
	}

	f.items = append(f.items, n)

	if n.moves == 0 {
		f.zeroAt = len(f.items) - 1
	}
}

func (f *File) Size() int {
	return len(f.items)
}

func (f *File) MixFilePart1(coordinates ...int) int64 {
	var (
		result int64
	)

	for i := 0; i < len(f.items); i++ {
		var (
			current = f.items[i]
			moves   = current.moves
		)

		if moves < 0 {
			f.mixBackwards(current, -moves)
		} else if moves > 0 {
			f.mixForward(current, moves)
		}
	}

	for i := 0; i < len(coordinates); i++ {
		current := f.items[f.zeroAt]

		for j := 0; j < coordinates[i]; j++ {
			current = current.next
		}

		result += current.moves
	}

	return result
}

// mixForward moves the value forward in the list.
func (f *File) mixForward(from *node, moves int64) {
	if moves == 0 {
		return
	}

	var (
		prev = from.prev
		next = from.next
	)

	prev.next = next
	next.prev = prev

	from.next = next.next
	from.prev = next

	f.mixForward(from.next, moves-1)
}

// mixBackwards moves the value backwards in the list.
func (f *File) mixBackwards(from *node, moves int64) {
	if moves == 0 {
		return
	}

	var (
		prev = from.prev
		next = from.next
	)

	prev.next = next
	next.prev = prev

	from.next = prev
	from.prev = prev.prev

	f.mixBackwards(from.prev, moves-1)
}
