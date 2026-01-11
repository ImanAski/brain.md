package command

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the brain repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		data, err := json.Marshal(map[string][]byte{
			"public":  pub,
			"private": priv,
		})
		if err != nil {
			return err
		}
		err = os.MkdirAll(GlobalContext.RootPath, 0700)
		if err != nil {
			return err
		}
		keysFile := filepath.Join(GlobalContext.RootPath, "keys.json")
		err = os.WriteFile(keysFile, data, 0600)
		if err != nil {
			return err
		}
		fmt.Println("keys.json created")
		return nil
	},
}

func init() {
	Register(initCmd)
}
