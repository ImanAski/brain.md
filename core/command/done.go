package command

import (
	"brain/core/object"
	"brain/core/types"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idHex := args[0]
		id, err := GlobalContext.Store.ResolveID(idHex)
		if err != nil {
			return err
		}

		o, err := GlobalContext.Store.Get(id)
		if err != nil {
			return err
		}
		if o == nil {
			return fmt.Errorf("object not found: %s", idHex)
		}

		if o.Type != "task" {
			return fmt.Errorf("object %s is not a task (type: %s)", idHex, o.Type)
		}

		var t types.Task
		if err := json.Unmarshal(o.Body, &t); err != nil {
			return fmt.Errorf("failed to unmarshal task: %w", err)
		}

		if t.Status == "done" {
			fmt.Printf("Task %s is already done.\n", idHex)
			return nil
		}

		t.Status = "done"
		body, _ := json.Marshal(t)

		newObj, err := object.New(
			GlobalContext.Keys["public"],
			"task",
			body,
			[]object.ID{o.ID},
			GlobalContext.Keys["private"],
		)
		if err != nil {
			return err
		}

		GlobalContext.Store.Put(newObj)
		fmt.Printf("Task %s marked as done (New state: %x)\n", idHex, newObj.ID)
		return nil
	},
}

func init() {
	Register(doneCmd)
}
