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
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homePlaceholder = "{{home}}"
var sshLocalFolder = homePlaceholder + "/.ssh"

func getProjectIDAndPath(providerConnectionString string) (string, string, error) {
	rgx := regexp.MustCompile("^.*:/(.*)@(.*)$")
	matches := rgx.FindAllStringSubmatch(providerConnectionString, 1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", "", errors.New("Cannot find project ID or Path in the provided string")
	}

	return matches[0][1], matches[0][2], nil
}

func replaceHomePlaceholder(path string) (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}
	rgx := regexp.MustCompile(homePlaceholder)
	return rgx.ReplaceAllString(path, homePath), nil
}

func createSSHConfigFile(path string, content []byte) error {
	file, err := os.Create(path)
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

func craftPath(path string) (string, error) {
	paths := strings.Split(path, "/")
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

		var wg sync.WaitGroup
		providerConnectionString := viper.GetString("provider")

		provider := providers.New(providerConnectionString)
		projectID, remotePath, err := getProjectIDAndPath(providerConnectionString)

		files, err := provider.GetFiles(projectID, remotePath)
		if err != nil {
			log.Fatal("Cannot get files from provider: " + err.Error())
		}

		for _, file := range files {

			wg.Add(1)
			fileID := file.ID
			filePath := file.Path

			go func() {
				defer wg.Done()
				fileContent, err := provider.GetFile(projectID, fileID)
				if err != nil {
					log.Fatal(err)
					return
				}

				filePath, err := craftPath(filePath)
				if err != nil {
					log.Fatal(err)
					return
				}

				err = createSSHConfigFile(filePath, fileContent)
				if err != nil {
					log.Fatal(err)
					return
				}

				log.Println("Created:", filePath)
			}()

		}

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
