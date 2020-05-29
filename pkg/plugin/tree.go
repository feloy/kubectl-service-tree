package plugin

import (
	"fmt"
	"io"
	"strings"
)

type Tree struct {
	Kind     string
	Name     string
	Info     *string
	Status   *bool
	Children []Tree
}

func NewTree(kind string, name string, info *string) Tree {
	return Tree{
		Kind:     kind,
		Name:     name,
		Info:     info,
		Children: []Tree{},
	}
}

func (o *Tree) setStatus(status bool) {
	o.Status = &status
}

func (o *Tree) addChild(child Tree) {
	o.Children = append(o.Children, child)
}

func (o Tree) display(w io.Writer, d ...int) {
	depth := 0
	if len(d) > 0 {
		depth = d[0]
	}

	status := ""
	if o.Status != nil {
		if *o.Status {
			status = ". "
		} else {
			status = "x "
		}
		depth--
	}

	indent := strings.Repeat(" ", depth)

	info := ""
	if o.Info != nil {
		info = fmt.Sprintf("(%s)", *o.Info)
	}

	fmt.Fprintf(w, "%s%s%s %s %s\n", indent, status, o.Kind, o.Name, info)
	for _, child := range o.Children {
		if o.Status != nil {
			depth += 2
		}
		child.display(w, depth+1)
	}
}
