package cmd

import (
	"hssh/messages"
	"hssh/services"
	"os"

	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:     "find",
	Aliases: []string{"f"},
	Short:   "Find host details using fzf",
	PreRun: func(cmd *cobra.Command, args []string) {
		services.Init(false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var host string
		if len(args) > 0 {
			host = args[0]
		}

		connection := services.Find(host)
		if connection.Name == "" {
			messages.NoConnection()
			os.Exit(0)
		}

		services.PrintConnectionDetails(&connection)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
