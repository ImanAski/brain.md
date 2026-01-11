package visualize

import (
	"brain/adapters/store/sqlite"
	"brain/core/object"
	"fmt"
)

type DotVisualizer struct{}

func init() {
	Register("dot", &DotVisualizer{})
}

func (d *DotVisualizer) Visualize(objects []*object.Object, links []sqlite.Link) error {
	fmt.Println("digraph G {")
	fmt.Println("  node [shape=box];")
	for _, o := range objects {
		fmt.Printf("  \"%x\" [label=\"%s: %x\"];\n", o.ID, o.Type, o.ID[:4])
	}
	for _, l := range links {
		fmt.Printf("  \"%x\" -> \"%x\";\n", l.Parent, l.Child)
	}
	fmt.Println("}")
	return nil
}
