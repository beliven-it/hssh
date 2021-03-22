package cmd

import (
	"hssh/controllers"
	"hssh/messages"
	"os"

	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:     "find",
	Aliases: []string{"f"},
	Short:   "Find host details using fzf",
	PreRun: func(cmd *cobra.Command, args []string) {
		controllers.Init(false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var host string
		if len(args) > 0 {
			host = args[0]
		}

		connection := controllers.Find(host)
		if connection.Name == "" {
			messages.NoConnection()
			os.Exit(0)
		}

		controllers.PrintConnectionDetails(&connection)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
