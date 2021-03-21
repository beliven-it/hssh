package cmd

import (
	"fmt"
	"hssh/controllers"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var syncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Short:   "Sync down hosts from the Git provider",
	PreRun: func(cmd *cobra.Command, args []string) {
		controllers.Init(2, false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("KK")
		controllers.Sync()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
