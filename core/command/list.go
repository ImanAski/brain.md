package command

import (
	"brain/core/handler"
	"fmt"
)

type ListCommand struct{}

func init() {
	Register("list", &ListCommand{})
}

func (c *ListCommand) Run(ctx *Context) error {
	for _, o := range ctx.Store.All() {
		h := handler.Get(o.Type)
		if h == nil {
			fmt.Printf("No handler for type %s\n", o.Type)
			continue
		}
		fmt.Printf("%x %s\n", o.ID, h.Render(o))
	}
	return nil
}
