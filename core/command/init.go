package command

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type InitCommand struct{}

func init() {
	Register("init", &InitCommand{})
}

func (c *InitCommand) Run(ctx *Context) error {
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
	err = os.MkdirAll(ctx.RootPath, 0700)
	if err != nil {
		return err
	}
	keysFile := filepath.Join(ctx.RootPath, "keys.json")
	err = os.WriteFile(keysFile, data, 0600)
	if err != nil {
		return err
	}
	fmt.Println("keys.json created")
	return nil
}
