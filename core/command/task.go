package command

import (
	"brain/core/object"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task <title> [parentID...]",
	Short: "Create a new task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		var parents []object.ID
		for i := 1; i < len(args); i++ {
			pID, err := hex.DecodeString(args[i])
			if err != nil {
				return fmt.Errorf("invalid parent ID: %s", args[i])
			}
			var id object.ID
			copy(id[:], pID)
			parents = append(parents, id)
		}

		due, _ := cmd.Flags().GetString("due")
		body, _ := json.Marshal(map[string]string{
			"title":  title,
			"status": "todo",
			"due":    due,
		})
		o, err := object.New(GlobalContext.Keys["public"], "task", body, parents, GlobalContext.Keys["private"])
		if err != nil {
			return err
		}
		fmt.Printf("Created task %x\n", o.ID)
		GlobalContext.Store.Put(o)
		return nil
	},
}

func init() {
	taskCmd.Flags().StringP("due", "d", "", "Due date (YYYY-MM-DD)")
	Register(taskCmd)
}
