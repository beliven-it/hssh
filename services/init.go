package services

import (
	"crypto/md5"
	"fmt"
	"hssh/config"
	"hssh/messages"
	"hssh/templates"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/spf13/viper"
)

const configInitializationStatusCode = 2

func upsertConfigSSH() error {
	file, err := os.OpenFile(config.SSHConfigFilePath, os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	oldContent, _ := os.ReadFile(config.SSHConfigFilePath)
	oldContentToString := string(oldContent)

	delimiterStart := "# HSSH start managed"
	delimiterEnd := "# HSSH end managed"
	includeString := "Include " + config.HSSHHostFolderName + "/*"

	row := delimiterStart + "\n" + includeString + "\n" + delimiterEnd + "\n\n"
	if isFilePathInConfigSSH(oldContentToString, row) {
		deleteRegex := regexp.MustCompile("(?ms)" + delimiterStart + ".*" + delimiterEnd + "\n\n")
		oldContentToString = deleteRegex.ReplaceAllString(oldContentToString, "")
	}

	_, err = file.WriteString(row + oldContentToString)
	if err != nil {
		return err
	}

	return nil
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
	upsertConfigSSH()
	cb(err)
}

func isHSSHInitialized() bool {
	h := md5.New()
	file, err := os.ReadFile(config.HSSHConfigFilePath)
	if err != nil {
		return false
	}

	return string(h.Sum(file)) != string(h.Sum([]byte(templates.Config)))
}

// Init ...
func Init(force bool) {
	isInit := isHSSHInitialized()
	if !force && !isInit {
		messages.ConfigNotEditedYet()
		os.Exit(0)
	}

	if !force && isInit {
		viper.SetConfigFile(config.HSSHConfigFilePath)
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			messages.ViperLoadError(err)
		}
		return
	}

	actions := []func(){
		func() {
			CreateHSSHConfig(func(err error, isNotConfigured bool) {
				messages.PrintStep(fmt.Sprintf("File %s", config.HSSHConfigFilePath), err)
				if isNotConfigured {
					messages.MustBeConfigured()
				}
			})
		},
		func() {
			CreateSSHConfig(func(err error) {
				messages.PrintStep(fmt.Sprintf("File %s", config.SSHConfigFilePath), err)
			})
		},
		func() {
			CreateHSSHHostFolder(func(err error) {
				messages.PrintStep(fmt.Sprintf("Folder %s", config.HSSHHostFolderPath), err)
				isEmpty, err := isFolderEmpty(config.HSSHHostFolderPath)
				if err != nil || isEmpty {
					Sync()
					messages.PrintStep("Automatic Sync", err)
				} else {
					messages.Print("Folder is not empty. The automatic sync will be ignored\n")
				}
			})
		},
	}

	for _, action := range actions {
		action()
	}
}
