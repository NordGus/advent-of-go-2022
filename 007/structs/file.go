package structs

import (
	"fmt"
	"strings"
)

type File struct {
	name string
	size uint64
}

func (f *File) ItemType() FilesystemItemType {
	return FileType
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() uint64 {
	return f.size
}

func (f *File) Print(level uint) {
	fmt.Printf("%v%v (%v)\n", strings.Repeat("\t", int(level)), f.name, f.size)
}
