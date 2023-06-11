package structs

import "sort"

type DirectorySize struct {
	Directory string
	Size      uint64
}

type bySize []DirectorySize

func (s bySize) Len() int {
	return len(s)
}

func (s bySize) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s bySize) Less(i, j int) bool {
	return s[i].Size < s[j].Size
}

func sortDirectorySize(s []DirectorySize) []DirectorySize {
	sort.Sort(bySize(s))

	return s
}
