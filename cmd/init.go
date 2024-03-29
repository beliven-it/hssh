package cmd

import (
	"hssh/controllers"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Init HSSH",
	Run: func(cmd *cobra.Command, args []string) {
		controllers.Init(true)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
