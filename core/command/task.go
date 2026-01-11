package command

import (
	"brain/core/object"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type TaskCommand struct{}

func init() {
	Register("task", &TaskCommand{})
}

func (c *TaskCommand) Run(ctx *Context) error {
	if len(ctx.Args) < 1 {
		return fmt.Errorf("usage: brain task <title> [parentID...]")
	}
	title := ctx.Args[0]
	var parents []object.ID
	for i := 1; i < len(ctx.Args); i++ {
		pID, err := hex.DecodeString(ctx.Args[i])
		if err != nil {
			return fmt.Errorf("invalid parent ID: %s", ctx.Args[i])
		}
		var id object.ID
		copy(id[:], pID)
		parents = append(parents, id)
	}

	body, _ := json.Marshal(map[string]string{"title": title, "status": "todo"})
	o, err := object.New(ctx.Keys["public"], "task", body, parents, ctx.Keys["private"])
	if err != nil {
		return err
	}
	fmt.Printf("Created task %x\n", o.ID)
	ctx.Store.Put(o)
	return nil
}
