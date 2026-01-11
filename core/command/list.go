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
		typ, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		}

		heads := GlobalContext.Store.Heads()
		for _, o := range heads {
			if typ != "" && o.Type != typ {
				continue
			}
			h := handler.Get(o.Type)
			if h == nil {
				fmt.Printf("%x (No handler for type %s)\n", o.ID, o.Type)
				continue
			}
			fmt.Printf("\033[33m%x\033[0m [\033[36m%s\033[0m] %s\n", o.ID[:4], o.Type, h.Render(o))
		}
		return nil
	},
}

func init() {
	listCmd.Flags().StringP("type", "t", "", "Filter by type")
	Register(listCmd)
}
