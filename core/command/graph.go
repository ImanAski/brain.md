package command

import (
	"brain/core/visualize"
	"fmt"
)

type GraphCommand struct{}

func init() {
	Register("graph", &GraphCommand{})
}

func (c *GraphCommand) Run(ctx *Context) error {
	vType := "dot"
	if len(ctx.Args) > 0 {
		vType = ctx.Args[0]
	}

	v := visualize.Get(vType)
	if v == nil {
		return fmt.Errorf("unknown visualizer: %s", vType)
	}

	objs := ctx.Store.All()
	links := ctx.Store.Links()

	return v.Visualize(objs, links)
}
