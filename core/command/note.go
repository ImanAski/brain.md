package command

import (
	"brain/core/object"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := args[0]
		body, _ := json.Marshal(map[string]string{"text": text})
		o, err := object.New(GlobalContext.Keys["public"], "note", body, nil, GlobalContext.Keys["private"])
		if err != nil {
			return err
		}
		fmt.Printf("Created note %x\n", o.ID)
		GlobalContext.Store.Put(o)
		return nil
	},
}

func init() {
	Register(noteCmd)
}
