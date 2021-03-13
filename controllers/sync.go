package controllers

import (
	"errors"
	"fmt"
	"hssh/providers"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

const homePlaceholder = "{{home}}"
const sshLocalFolder = homePlaceholder + "/.ssh"
const sshConfigFile = sshLocalFolder + "/config"

func getProjectIDAndPath(providerConnectionString string) (string, string, error) {
	rgx := regexp.MustCompile("^.*:/(.*)@(.*)$")
	matches := rgx.FindAllStringSubmatch(providerConnectionString, 1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", "", errors.New("Cannot find project ID or Path in the provided string")
	}

	return matches[0][1], matches[0][2], nil
}

func replaceHomePlaceholder(filePath string) (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return filePath, err
	}
	rgx := regexp.MustCompile(homePlaceholder)
	return rgx.ReplaceAllString(filePath, homePath), nil
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

func upsertConfigSSH(syncFile string) error {
	configFile, err := replaceHomePlaceholder(sshConfigFile)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(configFile, os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	oldContent, err := ioutil.ReadFile(configFile)
	oldContentToString := string(oldContent)

	var row = "# HSSH start managed\n" + "Include " + syncFile + "/*\n# HSSH end managed\n"
	if isFilePathInConfigSSH(oldContentToString, row) == true {
		deleteRegex := regexp.MustCompile("(?ms)# HSSH start managed.*# HSSH end managed\n")
		oldContentToString = deleteRegex.ReplaceAllString(oldContentToString, "")
	}

	_, err = file.WriteString(row + oldContentToString)
	if err != nil {
		return err
	}

	return nil
}

func craftPath(filePath string) (string, error) {
	paths := strings.Split(filePath, "/")
	deepFolder := len(paths) - 1
	folders := paths[0:deepFolder]
	foldersToCreate := ""

	if len(folders) > 0 {
		foldersToCreate = strings.Join(folders, "/")
	}

	filePath, err := replaceHomePlaceholder(sshLocalFolder)
	if err != nil {
		return "", err
	}

	if foldersToCreate != "" {
		filePath = filePath + "/" + foldersToCreate
		os.MkdirAll(filePath, os.ModePerm)
	}

	filePath = filePath + "/" + paths[deepFolder]

	return filePath, nil

}

// Sync ...
func Sync() {

	providerConnectionString := viper.GetString("provider")
	projectID, remotePath, err := getProjectIDAndPath(providerConnectionString)

	provider := providers.New(providerConnectionString)

	files, err := provider.GetFiles(projectID, remotePath)
	if err != nil {
		log.Fatal("Cannot get files from provider: " + err.Error())
	}

	// Create the entity in .ssh/config
	defer upsertConfigSSH(remotePath)

	for _, file := range files {

		fileID := file.ID
		filePath := file.Path

		func() {

			fileContent, err := provider.GetFile(projectID, fileID)
			if err != nil {
				log.Fatal(err)
				return
			}

			craftedPath, err := craftPath(filePath)
			if err != nil {
				log.Fatal(err)
				return
			}

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
