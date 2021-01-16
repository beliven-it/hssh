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
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"l"},
	Short: "List all available SSH aliases",
	Run: func(cmd *cobra.Command, args []string) {
		// Connection struct.
		type Connection struct {
			Name string
			Hostname string
			User string
			Port string
			IdentityFile string
		}

		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error finding home folder: %v\n", err)
			os.Exit(1)
		}

		// Get flags values
		colors, _:= cmd.Flags().GetBool("colors")

		files, err := ioutil.ReadDir(home + "/.ssh/" + viper.GetString("ssh_config_folder"))
		if err != nil {
			fmt.Printf("Error reading files in folder: %v\n", err)
			os.Exit(1)
		}

		content := ""
		for _, file := range files {
			data, err := ioutil.ReadFile(home + "/.ssh/" + viper.GetString("ssh_config_folder") + "/" + file.Name())
			if err != nil {
				fmt.Printf("File reading error: %v\n", err)
				os.Exit(1)
			}

			// Convert byte to string and add to content.
			content += string(data)
		}

		// Remove comments from content.
		content = regexp.MustCompile("(?m)^#.*").ReplaceAllString(content, "")

		// Remove empty lines from content.
		content = regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(content), "\n")

		// Split content into hosts.
		hosts := strings.Split(content, "Host ")

		// Map hosts into array of Connection struct.
		var connections []Connection
		for indexHost, host := range hosts {
			if indexHost == 0 {
				continue
			}

			host = strings.ReplaceAll(host, " ", "")

			var temp = Connection{}
			for indexParam, param := range strings.Split(host, "\n") {

				if indexParam == 0 {
					temp.Name = param
				} else {
					if strings.Contains(param, "Hostname") {
						temp.Hostname = strings.ReplaceAll(param, "Hostname", "")
					}

					if strings.Contains(param, "User") {
						temp.User = strings.ReplaceAll(param, "User", "")
					}

					if strings.Contains(param, "Port") {
						temp.Port = strings.ReplaceAll(param, "Port", "")
					}

					if strings.Contains(param, "IdentityFile") {
						temp.IdentityFile = strings.ReplaceAll(param, "IdentityFile", "")
					}
				}
			}

			connections = append(connections, temp)
		}

		// Sort alphabetically (case insensitive).
		sort.Slice(connections[:], func(i, j int) bool {
			return strings.ToLower(connections[i].Name) < strings.ToLower(connections[j].Name)
		})

		// Print connections.
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		for _, connection := range connections {
				if colors {
					fmt.Printf("%s -> %s@%s:%s %s\n", green(connection.Name), connection.User, connection.Hostname, red(connection.Port), yellow(connection.IdentityFile))
				} else {
					fmt.Printf("%s -> %s@%s:%s %s\n", connection.Name, connection.User, connection.Hostname, connection.Port, connection.IdentityFile)
				}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("colors", "c", false, "List hosts with color highlights.")
}
