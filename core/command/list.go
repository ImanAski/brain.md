package command

import (
	"brain/core/handler"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all items in the brain",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, o := range GlobalContext.Store.All() {
			h := handler.Get(o.Type)
			if h == nil {
				fmt.Printf("No handler for type %s\n", o.Type)
				continue
			}
			fmt.Printf("%s\n", h.Render(o))
		}
		return nil
	},
}

func init() {
	Register(listCmd)
}
