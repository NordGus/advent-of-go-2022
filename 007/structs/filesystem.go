package structs

import "sort"

type FilesystemItemType uint8

const (
	DirectoryType FilesystemItemType = iota
	FileType
)

type FilesystemItem interface {
	ItemType() FilesystemItemType
	Name() string
	Size() uint64
	Print(level uint)
}

type Filesystem struct {
	root             *Directory
	currentDirectory *Directory
}

func NewFilesystem() Filesystem {
	root := Directory{
		name: "/",
	}

	return Filesystem{
		root:             &root,
		currentDirectory: &root,
	}
}

func (fs *Filesystem) Stream(name string, meta string) {
	err := fs.currentDirectory.addItem(name, meta)
	if err != nil {
		panic(err)
	}
}

func (fs *Filesystem) CD(target string) {
	oldDirectory := fs.currentDirectory

	if target == fs.root.name && target == fs.currentDirectory.name {
		return
	}

	newDirectory, err := oldDirectory.navigateTo(target)
	if err != nil {
		panic(err)
	}

	fs.currentDirectory = newDirectory
}

func (fs *Filesystem) DirectorySizes() []DirectorySize {
	dirs := fs.root.getDirectoriesDirectories()
	results := make([]DirectorySize, len(dirs))

	for i := 0; i < len(dirs); i++ {
		results[i] = DirectorySize{Directory: dirs[i].Path(), Size: dirs[i].Size()}
	}

	sort.Sort(bySize(results))

	return results
}

func (fs *Filesystem) PrintFromRoot() {
	fs.root.Print(0)
}
