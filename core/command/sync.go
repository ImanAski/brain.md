package command

import (
	"brain/core/object"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize the brain with other peers",
}

var exportCmd = &cobra.Command{
	Use:   "export <file>",
	Short: "Export the entire brain to a bundle file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		objects := GlobalContext.Store.All()

		data, err := json.MarshalIndent(objects, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return err
		}

		fmt.Printf("Exported %d objects to %s\n", len(objects), filePath)
		return nil
	},
}

var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import objects from a bundle file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		var objects []*object.Object
		if err := json.Unmarshal(data, &objects); err != nil {
			return err
		}

		count := 0
		for _, o := range objects {
			// In a real system, we would verify signatures here before Put.
			// o.Verify() exists in core/object/object.go
			if !o.Verify() {
				fmt.Printf("Warning: Skipping invalid object %x\n", o.ID)
				continue
			}
			GlobalContext.Store.Put(o)
			count++
		}

		fmt.Printf("Imported %d objects from %s\n", count, filePath)
		return nil
	},
}

func init() {
	syncCmd.AddCommand(exportCmd)
	syncCmd.AddCommand(importCmd)
	Register(syncCmd)
}
