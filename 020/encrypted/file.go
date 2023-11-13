package encrypted

const (
	defaultFileSize = 5000
)

type File struct {
	original []int64
}

func New() File {
	return File{
		original: make([]int64, 0, defaultFileSize),
	}
}

func (f *File) AddItem(item int64) {
	f.original = append(f.original, item)
}

func (f *File) Size() int {
	return len(f.original)
}
