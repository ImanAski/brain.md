package command

import (
	"brain/adapters/store/sqlite"
)

type Context struct {
	Args     []string
	RootPath string
	Store    *sqlite.Store
	Keys     map[string][]byte
}

type Command interface {
	Run(ctx *Context) error
}

var registry = make(map[string]Command)

func Register(name string, cmd Command) {
	registry[name] = cmd
}

func Get(name string) Command {
	return registry[name]
}
