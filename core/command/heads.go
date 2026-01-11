package command

import (
	"brain/core/handler"
	"fmt"

	"github.com/spf13/cobra"
)

var headsCmd = &cobra.Command{
	Use:   "heads",
	Short: "List all head objects in the brain",
	RunE: func(cmd *cobra.Command, args []string) error {
		heads := GlobalContext.Store.Heads()
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
	},
}

func init() {
	Register(headsCmd)
}
