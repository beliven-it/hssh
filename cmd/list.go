package cmd

import (
	"hssh/services"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all available hosts",
	PreRun: func(cmd *cobra.Command, args []string) {
		services.Init(false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		colors, _ := cmd.Flags().GetBool("colors")
		connections := services.List()

		for _, connection := range connections {
			services.PrintConnection(&connection, colors)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("colors", "c", false, "List hosts with color highlights.")
}
