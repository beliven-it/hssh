package cmd

import (
	"hssh/services"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Short:   "Sync down hosts from the Git provider",
	PreRun: func(cmd *cobra.Command, args []string) {
		services.Init(false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		services.Sync()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
