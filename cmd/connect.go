package cmd

import (
	"hssh/services"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:     "connect",
	Aliases: []string{"c"},
	Short:   "Search and connect to host using fzf",
	PreRun: func(cmd *cobra.Command, args []string) {
		services.Init(false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var host string
		if len(args) > 0 {
			host = args[0]
		}

		services.Connect(host)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
