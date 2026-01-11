package command

import (
	"brain/adapters/store/sqlite"

	"github.com/spf13/cobra"
)

type Context struct {
	RootPath string
	Store    *sqlite.Store
	Keys     map[string][]byte
}

var GlobalContext = &Context{}

var RootCmd = &cobra.Command{
	Use:   "brain",
	Short: "Brain is a tool for managing your personal knowledge graph",
}

func Register(cmd *cobra.Command) {
	RootCmd.AddCommand(cmd)
}
