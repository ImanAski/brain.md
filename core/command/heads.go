package command

import (
	"brain/core/handler"
	"fmt"
)

type HeadsCommand struct{}

func init() {
	Register("heads", &HeadsCommand{})
}

func (c *HeadsCommand) Run(ctx *Context) error {
	heads := ctx.Store.Heads()
	if len(heads) == 0 {
		fmt.Println("No heads found.")
		return nil
	}

	for _, o := range heads {
		h := handler.Get(o.Type)
		if h == nil {
			fmt.Printf("%x (No handler for type %s)\n", o.ID, o.Type)
			continue
		}
		fmt.Printf("%x [%s] %s\n", o.ID, o.Type, h.Render(o))
	}
	return nil
}
