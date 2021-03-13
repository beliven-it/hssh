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

func printConnection(connection *models.Connection) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgHiMagenta).SprintFunc()

	fmt.Printf("\nName: %s\nHostname: %s\nUser: %s\nPort: %s\nIdentity: %s\n",
		green(connection.Name), magenta(connection.Hostname), blue(connection.User), red(connection.Port), yellow(connection.IdentityFile))
}

// listCmd represents the list command
var findCmd = &cobra.Command{
	Use:     "find",
	Aliases: []string{"f"},
	Short:   "Find the details about a connection",
	Run: func(cmd *cobra.Command, args []string) {
		connection := controllers.Find()
		printConnection(&connection)

	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
