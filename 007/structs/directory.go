package structs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Directory struct {
	parent   *Directory
	children []FilesystemItem
	name     string
}

func (d *Directory) addItem(name string, meta string) error {
	var itemType FilesystemItemType

	switch meta {
	case "dir":
		itemType = DirectoryType
	default:
		itemType = FileType
	}

	for _, child := range d.children {
		if child.Name() == name && child.ItemType() == itemType {
			switch itemType {
			case DirectoryType:
				return fmt.Errorf("%v directory already exists", name)
			default:
				return fmt.Errorf("%v file already exists", name)
			}
		}
	}

	switch itemType {
	case DirectoryType:
		dir := Directory{name: name, parent: d}
		d.children = append(d.children, &dir)
	default:
		size, err := strconv.ParseUint(meta, 10, 0)
		if err != nil {
			return err
		}

		file := File{name: name, size: size}
		d.children = append(d.children, &file)
	}

	return nil
}

func (d *Directory) navigateTo(target string) (*Directory, error) {
	if target == ".." && d.parent == nil {
		return d.parent, errors.New("root directory")
	}

	if target == ".." && d.parent != nil {
		return d.parent, nil
	}

	for _, child := range d.children {
		if child.Name() == target && child.ItemType() == DirectoryType {
			return child.(*Directory), nil
		}
	}

	return nil, fmt.Errorf("directory (%v) not found in %v", target, d.Path())
}

func (d *Directory) Path() string {
	if d.parent == nil {
		return d.name
	}

	return fmt.Sprintf("%v%v/", d.parent.Path(), d.name)
}

func (d *Directory) ItemType() FilesystemItemType {
	return DirectoryType
}

func (d *Directory) Name() string {
	return d.name
}

func (d *Directory) Size() uint64 {
	var size uint64 = 0

	for _, child := range d.children {
		size += child.Size()
	}

	return size
}

func (d *Directory) Print(level uint) {
	fmt.Printf("%v- %v (dir)\n", strings.Repeat("\t", int(level)), d.name)

	for _, child := range d.children {
		child.Print(level + 1)
	}
}

func (d *Directory) hasChildDirectories() bool {
	for _, child := range d.children {
		if child.ItemType() == DirectoryType {
			return true
		}
	}

	return false
}

func (d *Directory) getDirectoriesDirectories() []Directory {
	if !d.hasChildDirectories() {
		return []Directory{*d}
	}

	dirs := make([]Directory, 0)

	for _, child := range d.children {
		if child.ItemType() == DirectoryType {
			dir := child.(*Directory)
			dirs = append(dirs, dir.getDirectoriesDirectories()...)
		}
	}

	dirs = append(dirs, *d)

	return dirs
}
