package controllers

import (
	"errors"
	"fmt"
	"hssh/config"
	"hssh/providers"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

func getProjectIDAndPath(providerConnectionString string) (string, string, error) {
	rgx := regexp.MustCompile("^.*:/(.*)@(.*)$")
	matches := rgx.FindAllStringSubmatch(providerConnectionString, 1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", "", errors.New("Cannot find project ID or Path in the provided string")
	}

	return matches[0][1], matches[0][2], nil
}

func createSSHConfigFile(filePath string, content []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}

func isFilePathInConfigSSH(content string, filePath string) bool {
	replacer := regexp.MustCompile("\\*")

	rgx := regexp.MustCompile("(?m)" + replacer.ReplaceAllString(filePath, "\\*"))
	return rgx.MatchString(content)
}

func upsertConfigSSH() error {
	file, err := os.OpenFile(config.SSHConfigFilePath, os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	oldContent, err := ioutil.ReadFile(config.SSHConfigFilePath)
	oldContentToString := string(oldContent)

	delimiterStart := "# HSSH start managed"
	delimiterEnd := "# HSSH end managed"
	includeString := "Include " + config.HSSHHostFolderName + "/*"

	var row = delimiterStart + "\n" + includeString + "\n" + delimiterEnd + "\n"
	if isFilePathInConfigSSH(oldContentToString, row) == true {
		deleteRegex := regexp.MustCompile("(?ms)" + delimiterStart + ".*" + delimiterEnd + "\n")
		oldContentToString = deleteRegex.ReplaceAllString(oldContentToString, "")
	}

	_, err = file.WriteString(row + oldContentToString)
	if err != nil {
		return err
	}

	return nil
}

func craftPath(filePath string) string {
	paths := strings.Split(filePath, "/")
	fileName := paths[len(paths)-1]

	filePath = config.HSSHHostFolderPath + "/" + fileName

	return filePath
}

// Sync ...
func Sync() {
	providerConnectionString := viper.GetString("provider")
	projectID, remotePath, err := getProjectIDAndPath(providerConnectionString)

	provider := providers.New(providerConnectionString)

	files, err := provider.GetFiles(projectID, remotePath)
	if err != nil {
		fmt.Println("Cannot get files from provider: " + err.Error())
		os.Exit(1)
	}

	// Create the entity in .ssh/config
	defer upsertConfigSSH()

	for _, file := range files {

		fileID := file.ID
		filePath := file.Path

		func() {
			fileContent, err := provider.GetFile(projectID, fileID)
			if err != nil {
				log.Fatal(err)
				return
			}

			craftedPath := craftPath(filePath)
			err = createSSHConfigFile(craftedPath, fileContent)
			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Println("Create file", craftedPath)

			return
		}()

	}
}
