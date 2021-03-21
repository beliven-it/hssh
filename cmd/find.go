package cmd

import (
	"fmt"
	"hssh/controllers"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
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
			fmt.Println("Cannot find the details of the host.\nRetry using fuzzysearch:\n\nhssh find")
			os.Exit(0)
		}

		controllers.PrintConnectionDetails(&connection)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
