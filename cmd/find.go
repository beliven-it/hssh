package cmd

import (
	"hssh/controllers"

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
		connection := controllers.Find()
		controllers.PrintConnectionDetails(&connection)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
