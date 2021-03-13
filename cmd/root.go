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
	"hssh/templates"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const configInitializationStatusCode = 2

// Version of the app provided
// in build phase
var Version string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     "hssh",
	Short:   "Heply CLI to simplify the management of SSH hosts",
	Version: Version,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initSSHConfig check if file ~/.ssh/config
// exist and create it if not
func initRequiredHomeSpaceFile(filePath string, template string) (string, int, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set config file path
	var configPath string = home + filePath

	// Create needed folders if not exist
	err = os.MkdirAll(path.Dir(configPath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating folders: %v\n", err)
		return "", 1, err
	}

	// Create config file starting from template if not exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, err := os.Create(configPath)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return "", 1, err
		}

		defer file.Close()
		file.WriteString(template)

		fmt.Printf("Created missing %v file!", configPath)
		return "", configInitializationStatusCode, nil
	}

	return configPath, 0, nil

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Check or create configuration file (config.yml)
	configPath, statusCode, err := initRequiredHomeSpaceFile("/.config/hssh/config.yml", templates.Config)
	if err != nil {
		fmt.Println("An error occured during config.yml initialization")
		os.Exit(1)
	}

	if statusCode == configInitializationStatusCode {
		fmt.Println("Before starting to use hssh edit the newly created configuration file")
		os.Exit(1)
	}

	// Check or create configuration ssh file (.ssh/config)
	// If not exist the file will created empty
	_, _, err = initRequiredHomeSpaceFile("/.ssh/config", "")
	if err != nil {
		fmt.Println("An error occured during ssh config initialization")
		os.Exit(1)
	}

	// Search "config.yml" file in "$HOME/.config/hssh" directory.
	viper.SetConfigFile(configPath)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}
}
