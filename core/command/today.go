package command

import (
	"brain/core/handler"
	"brain/core/types"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show tasks due today or overdue",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now().Format("2006-01-02")
		heads := GlobalContext.Store.Heads()

		found := false
		for _, o := range heads {
			if o.Type != "task" {
				continue
			}

			var t types.Task
			if err := json.Unmarshal(o.Body, &t); err != nil {
				continue
			}

			if t.Status == "done" {
				continue
			}

			// If due date is empty, we don't show it in "today" unless we want "today" to be all pending tasks.
			// overview.md says "due <= today".
			if t.Due != "" && t.Due <= now {
				h := handler.Get(o.Type)
				if h != nil {
					fmt.Printf("\033[33m%x\033[0m %s\n", o.ID[:4], h.Render(o))
					if t.Due < now {
						fmt.Printf("   \033[31m(Overdue: %s)\033[0m\n", t.Due)
					}
					found = true
				}
			}
		}

		if !found {
			fmt.Println("No tasks for today!")
		}
		return nil
	},
}

func init() {
	Register(todayCmd)
}
