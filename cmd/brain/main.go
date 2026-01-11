package main

import (
	"brain/adapters/store/sqlite"
	"brain/core/command"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_ "brain/core/command"
	_ "brain/core/types"
	_ "brain/core/visualize"
)

const VERSION = "0.0.1"

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	rootPath, err := generateRootPath()
	if err != nil {
		panic(err)
	}

	command.GlobalContext.RootPath = rootPath
	command.RootCmd.Version = VERSION

	command.RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "init" || cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Parent() == nil {
			return nil
		}

		keysFile := filepath.Join(command.GlobalContext.RootPath, "keys.json")
		dbFile := filepath.Join(command.GlobalContext.RootPath, "index.db")

		keyData, err := os.ReadFile(keysFile)
		if err != nil {
			return fmt.Errorf("keys.json not found. Run 'brain init' first")
		}
		var keys map[string][]byte
		if err = json.Unmarshal(keyData, &keys); err != nil {
			return err
		}
		command.GlobalContext.Keys = keys
		command.GlobalContext.Store = sqlite.Open(dbFile)
		return nil
	}

	if err := command.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func generateRootPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	rootPath := filepath.Join(home, ".brain")
	return rootPath, nil
}
