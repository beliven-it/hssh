package cmd

import (
	"hssh/services"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Init HSSH",
	Run: func(cmd *cobra.Command, args []string) {
		services.Init(true)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
