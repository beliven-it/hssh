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
	"errors"
	"hssh/providers"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homePlaceholder = "{{home}}"
var sshLocalFolder = homePlaceholder + "/.ssh"
var sshConfigFile = sshLocalFolder + "/config"

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
	rgx := regexp.MustCompile(filePath)
	return rgx.MatchString(content)
}

func addFilePathToConfigSSH(syncFile string) error {
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

	if isFilePathInConfigSSH(oldContentToString, syncFile) == true {
		return nil
	}

	var row = "# File inserted with HSSH \n" + syncFile + "\n" + oldContentToString

	_, err = file.WriteString(row)
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

// listCmd represents the list command
var syncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Short:   "Sync connection files from the provider selected",
	Run: func(cmd *cobra.Command, args []string) {

		providerConnectionString := viper.GetString("provider")

		provider := providers.New(providerConnectionString)
		projectID, remotePath, err := getProjectIDAndPath(providerConnectionString)

		files, err := provider.GetFiles(projectID, remotePath)
		if err != nil {
			log.Fatal("Cannot get files from provider: " + err.Error())
		}

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

				err = addFilePathToConfigSSH(craftedPath)
				if err != nil {
					log.Fatal(err)
					return
				}

				return
			}()

		}

	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
