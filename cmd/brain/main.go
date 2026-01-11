package main

import (
	"brain/adapters/store/sqlite"
	"brain/core/command"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "brain/core/command"
	_ "brain/core/types"
	_ "brain/core/visualize"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: brain <command>")
		os.Exit(1)
	}

	rootPath, err := generateRootPath()
	if err != nil {
		panic(err)
	}

	cmdName := os.Args[1]
	ctx := &command.Context{
		Args:     os.Args[2:],
		RootPath: rootPath,
	}

	if cmdName != "init" {
		keysFile := filepath.Join(rootPath, "keys.json")
		dbFile := filepath.Join(rootPath, "index.db")

		keyData, err := os.ReadFile(keysFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: keys.json not found. Run 'brain init' first.\n")
			os.Exit(1)
		}
		var keys map[string][]byte
		if err = json.Unmarshal(keyData, &keys); err != nil {
			panic(err)
		}
		ctx.Keys = keys
		ctx.Store = sqlite.Open(dbFile)
	}

	c := command.Get(cmdName)
	if c == nil {
		fmt.Printf("Unknown command: %s\n", cmdName)
		os.Exit(1)
	}

	if err := c.Run(ctx); err != nil {
		panic(err)
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
