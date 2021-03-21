package cmd

import (
	"hssh/controllers"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all available hosts",
	PreRun: func(cmd *cobra.Command, args []string) {
		controllers.Init(-1, false)
	},
	Run: func(cmd *cobra.Command, args []string) {
		colors, _ := cmd.Flags().GetBool("colors")
		connections := controllers.List()

		for _, connection := range connections {
			controllers.PrintConnection(&connection, colors)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("colors", "c", false, "List hosts with color highlights.")
}
