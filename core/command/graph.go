package command

import (
	"brain/core/visualize"
	"fmt"

	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph [type]",
	Short: "Visualize the brain graph",
	RunE: func(cmd *cobra.Command, args []string) error {
		vType := "dot"
		if len(args) > 0 {
			vType = args[0]
		}

		v := visualize.Get(vType)
		if v == nil {
			return fmt.Errorf("unknown visualizer: %s", vType)
		}

		objs := GlobalContext.Store.All()
		links := GlobalContext.Store.Links()

		return v.Visualize(objs, links)
	},
}

func init() {
	Register(graphCmd)
}
