package controllers

import (
	"fmt"
	"hssh/config"
	"hssh/templates"
	"io"
	"os"
	"path"

	"github.com/spf13/viper"
)

const configInitializationStatusCode = 2

func getStatusByError(err error) string {
	if err != nil {
		return "NOK"
	}
	return "OK"
}

// InitSSHConfig check if file ~/.ssh/config
// exist and create it if not
func initRequiredHomeSpaceFile(configPath string, template string) (int, error) {
	// Create needed folders if not exist
	err := os.MkdirAll(path.Dir(configPath), os.ModePerm)
	if err != nil {
		return 1, err
	}

	// Create config file starting from template if not exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, err := os.Create(configPath)
		if err != nil {
			return 1, err
		}

		defer file.Close()
		file.WriteString(template)

		return configInitializationStatusCode, nil
	}

	return 0, nil
}

func isFolderEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// CreateHSSHHostFolder ...
func CreateHSSHHostFolder(cb func(error)) {
	err := os.MkdirAll(config.HSSHHostFolderPath, os.ModePerm)
	cb(err)
}

// ExecuteFirstSync ...
func ExecuteFirstSync(cb func(error)) {
	isEmpty, err := isFolderEmpty(config.HSSHHostFolderPath)
	if err != nil || isEmpty == true {
		Sync()
	}

	cb(err)
}

// CreateHSSHConfig ...
// Check or create configuration file (config.yml)
func CreateHSSHConfig(cb func(error, bool)) {
	statusCode, err := initRequiredHomeSpaceFile(config.HSSHConfigFilePath, templates.Config)
	if err != nil {
		os.Exit(1)
	}

	cb(err, statusCode == configInitializationStatusCode)
}

// CreateSSHConfig ...
// Check or create configuration ssh file (.ssh/config)
// If not exist the file will created empty
func CreateSSHConfig(cb func(error)) {
	_, err := initRequiredHomeSpaceFile(config.SSHConfigFilePath, "")
	cb(err)
}

// Init ...
func Init(steps int, verbose bool) {
	actions := []func(){
		func() {
			CreateHSSHConfig(func(err error, isNotConfigured bool) {
				viper.SetConfigFile(config.HSSHConfigFilePath)
				viper.AutomaticEnv()
				if err := viper.ReadInConfig(); err != nil {
					fmt.Printf("Error reading config file: %v.\n", err)
					os.Exit(1)
				}
				if verbose {
					fmt.Printf("[%s] File %s.\n", getStatusByError(err), config.HSSHConfigFilePath)
				}
				if err != nil {
					fmt.Printf("%s\n", err.Error())
					os.Exit(1)
				}
				if isNotConfigured == true {
					fmt.Printf("NOTE! The file must be configured before using the application.\n")
				}

			})
		},
		func() {
			CreateSSHConfig(func(err error) {
				if verbose {
					fmt.Printf("[%s] File %s.\n", getStatusByError(err), config.SSHConfigFilePath)
				}
				if err != nil {
					fmt.Printf("%s\n", err.Error())
					os.Exit(1)
				}
			})
		},
		func() {
			CreateHSSHHostFolder(func(err error) {
				if verbose {
					fmt.Printf("[%s] Folder %s.\n", getStatusByError(err), config.HSSHHostFolderPath)
				}
				if err != nil {
					fmt.Printf("%s\n", err.Error())
					os.Exit(1)
				}
			})
		},
		func() {
			ExecuteFirstSync(func(err error) {
				if verbose {
					fmt.Printf("[%s] First Sync.\n", getStatusByError(err))
				}
				if err != nil {
					fmt.Printf("%s\n", err.Error())
					os.Exit(1)
				}
			})
		},
	}

	if steps >= 0 {
		actions = actions[0 : steps+1]
	}

	for _, action := range actions {
		action()
	}
}
