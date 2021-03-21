package controllers

import (
	"errors"
	"fmt"
	"hssh/config"
	"hssh/models"
	"hssh/providers"
	"os"
	"regexp"
	"strings"
	"sync"

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

func isFilePathInConfigSSH(content string, filePath string) bool {
	replacer := regexp.MustCompile("\\*")

	rgx := regexp.MustCompile("(?m)" + replacer.ReplaceAllString(filePath, "\\*"))
	return rgx.MatchString(content)
}

func craftPath(filePath string) string {
	paths := strings.Split(filePath, "/")
	fileName := paths[len(paths)-1]

	filePath = config.HSSHHostFolderPath + "/" + fileName

	return filePath
}

func syncWithProvider(providerConnection string) {
	projectID, remotePath, err := getProjectIDAndPath(providerConnection)

	provider := providers.New(providerConnection)

	files, err := provider.GetFiles(projectID, remotePath)
	if err != nil {
		fmt.Println("Cannot get files from provider: " + err.Error())
		os.Exit(1)
	}

	var wg = new(sync.WaitGroup)

	for _, file := range files {

		fileID := file.ID
		filePath := file.Path

		wg.Add(1)
		go func(path string, id string) {
			defer wg.Done()

			fileContent, err := provider.GetFile(projectID, id)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			host := models.NewHost(craftPath(path) + "." + provider.GetDriver())

			err = host.Create(fileContent)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println("File created:", host.GetPath())
		}(filePath, fileID)

	}
	wg.Wait()
}

// Sync ...
func Sync() {
	singleProvider := viper.GetString("provider")
	multiProvider := viper.GetStringSlice("providers")

	multiProvider = append(multiProvider, singleProvider)
	multiProvider = unique(multiProvider)

	var wg = new(sync.WaitGroup)

	for _, provider := range multiProvider {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if p == "" {
				return
			}

			syncWithProvider(p)
		}(provider)
	}

	wg.Wait()

}
