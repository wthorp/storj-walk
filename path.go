package main

import (
	"strings"

	"github.com/lxn/walk"
)

type Path struct {
	name     string
	parent   *Path
	children []*Path
}

var folderIcon *walk.Icon

func init() {
	folderIcon, _ = walk.NewIconFromSysDLL("SHELL32", 4)
}

func NewPath(name string, parent *Path) *Path {
	return &Path{name: name, parent: parent, children: []*Path{}}
}

var _ walk.TreeItem = new(Path)

func (d *Path) Text() string {
	return d.name
}

func (d *Path) Parent() walk.TreeItem {
	if d.parent == nil {
		// We can't simply return d.parent in this case, because the interface
		// value then would not be nil.
		return nil
	}
	return d.parent
}

func (d *Path) ChildCount() int {
	return len(d.children)
}

func (d *Path) ChildAt(index int) walk.TreeItem {
	return d.children[index]
}

func (d *Path) Image() interface{} {
	return folderIcon
}

func (d *Path) ResetChildren() error {
	return nil
}

func (d *Path) Path() (bucket, path string) {
	elems := []string{d.name}
	dir, _ := d.Parent().(*Path)
	for dir != nil {
		elems = append([]string{dir.name}, elems...)
		dir, _ = dir.Parent().(*Path)
	}
	return elems[0], strings.Join(elems[1:], "/")
}
