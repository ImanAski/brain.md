package command

import (
	"brain/core/handler"
	"brain/core/object"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show object graph history",
	RunE: func(cmd *cobra.Command, args []string) error {
		graph, _ := cmd.Flags().GetBool("graph")
		oneline, _ := cmd.Flags().GetBool("oneline")

		objects := GlobalContext.Store.All()
		// Simple topological sort or just reverse chronological for basic log
		// Since we don't have timestamps, we can sort by hash or just use list order (insertion order in sqlite)

		if graph {
			renderLogGraph(objects, oneline)
		} else {
			renderLogPlain(objects, oneline)
		}
		return nil
	},
}

func renderLogPlain(objects []*object.Object, oneline bool) {
	for i := len(objects) - 1; i >= 0; i-- {
		o := objects[i]
		h := handler.Get(o.Type)
		renderText := ""
		if h != nil {
			renderText = h.Render(o)
		}

		if oneline {
			fmt.Printf("%x [%s] %s\n", o.ID[:4], o.Type, renderText)
		} else {
			fmt.Printf("Object: %x\n", o.ID)
			fmt.Printf("Type:   %s\n", o.Type)
			if len(o.Parents) > 0 {
				parentSHAs := []string{}
				for _, p := range o.Parents {
					parentSHAs = append(parentSHAs, fmt.Sprintf("%x", p[:4]))
				}
				fmt.Printf("Parents: %s\n", strings.Join(parentSHAs, ", "))
			}
			fmt.Printf("Body:    %s\n", renderText)
			fmt.Println("---")
		}
	}
}

func renderLogGraph(objects []*object.Object, oneline bool) {
	// A very basic "git log --graph" style ASCII renderer
	// For now, let's just do a simple indentation-based or prefix-based view
	// since a true ASCII DAG renderer is complex.

	// Map to find objects by ID
	objMap := make(map[object.ID]*object.Object)
	for _, o := range objects {
		objMap[o.ID] = o
	}

	// We'll iterate through heads and trace back?
	// Actually, let's just do a linear log with parent markers for now to keep it simple but informative.
	for i := len(objects) - 1; i >= 0; i-- {
		o := objects[i]
		h := handler.Get(o.Type)
		renderText := ""
		if h != nil {
			renderText = h.Render(o)
		}

		prefix := "* "
		if len(o.Parents) > 1 {
			prefix = "M " // Merge
		}

		fmt.Printf("%s %x [%s] %s\n", prefix, o.ID[:4], o.Type, renderText)
		for _, p := range o.Parents {
			fmt.Printf("|  -> %x\n", p[:4])
		}
	}
}

func init() {
	logCmd.Flags().BoolP("graph", "g", false, "Show graph visualization")
	logCmd.Flags().Bool("oneline", false, "Show compact output")
	Register(logCmd)
}
