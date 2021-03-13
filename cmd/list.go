/*
Copyright © 2020 Heply SRL <hello@heply.it>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"hssh/controllers"
	"hssh/models"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func printConnections(connections []models.Connection, withColor bool) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	for _, connection := range connections {
		if withColor {
			fmt.Printf("%s -> %s@%s:%s %s\n", green(connection.Name), connection.User, connection.Hostname, red(connection.Port), yellow(connection.IdentityFile))
		} else {
			fmt.Printf("%s -> %s@%s:%s %s\n", connection.Name, connection.User, connection.Hostname, connection.Port, connection.IdentityFile)
		}
	}
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all available hosts",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags values
		colors, _ := cmd.Flags().GetBool("colors")

		connections := controllers.List()

		printConnections(connections, colors)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("colors", "c", false, "List hosts with color highlights.")
}
