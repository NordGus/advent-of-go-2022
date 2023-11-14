package encrypted

const (
	defaultFileSize = 5000
)

type File struct {
	original []int64
	mixed    []int
}

func New() File {
	return File{
		original: make([]int64, 0, defaultFileSize),
		mixed:    make([]int, 0, defaultFileSize),
	}
}

func (f *File) AddItem(item int64) {
	f.original = append(f.original, item)
	f.mixed = append(f.mixed, len(f.original)-1)
}

func (f *File) Size() int {
	return len(f.original)
}

func (f *File) MixFilePart1(coordinates ...int) int64 {
	var (
		result int64
	)

	for ogIndex := 0; ogIndex < len(f.original); ogIndex++ {
		var (
			from  = f.toMixedIndex(ogIndex)
			to    = f.getTargetMixedIndex(int(f.original[ogIndex]), from)
			value = f.mixed[from]
		)

		f.mix(value, from, to)
	}

	for i := 0; i < len(f.original); i++ {
		if f.original[i] == 0 {
			index := f.toMixedIndex(i)

			for i := 0; i < len(coordinates); i++ {
				result += f.original[f.mixed[(index+coordinates[i])%len(f.mixed)]]
			}

			break
		}
	}

	return result
}

// toMixedIndex is a helper function that takes the index of an item in the original slice and translates it to one of the mixed slice.
// It returns -1 if it can't find the item.
func (f *File) toMixedIndex(index int) int {
	for i := 0; i < len(f.mixed); i++ {
		if f.mixed[i] == index {
			return i
		}
	}

	return -1
}

// getTargetMixedIndex is a helper function that takes the moves dictated by the original slice values and translates it to an index of the mixed slice.
func (f *File) getTargetMixedIndex(moves int, from int) int {
	var target int

	if moves == 0 {
		return from
	}

	if moves < 0 {
		target = (moves + from - 1) % len(f.mixed)
	} else if moves+from > len(f.mixed) {
		target = (moves + from + 1) % len(f.mixed)
	} else {
		target = (moves + from) % len(f.mixed)
	}

	if target < 0 {
		return target + len(f.mixed)
	}

	return target
}

func (f *File) mix(value int, from int, to int) {
	f.mixed = append(f.mixed[:from], f.mixed[from+1:]...) // remove value

	var (
		first  = make([]int, len(f.mixed[:to]), len(f.mixed)+1)
		second = make([]int, len(f.mixed[to:]))
	)

	copy(first, f.mixed[:to])
	copy(second, f.mixed[to:])

	first = append(first, value)

	f.mixed = append(first, second...)
}
