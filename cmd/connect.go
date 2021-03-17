package cmd

import (
	"hssh/controllers"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:     "connect",
	Aliases: []string{"c"},
	Short:   "Search and connect to host using fzf",
	PreRun: func(cmd *cobra.Command, args []string) {
		controllers.Init(false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var host string
		if len(args) > 0 {
			host = args[0]
		}

		controllers.Connect(host)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
