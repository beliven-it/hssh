package controllers

import (
	"errors"
	"fmt"
	"hssh/config"
	"hssh/messages"
	"hssh/models"
	"hssh/providers"
	"os"
	"path/filepath"
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

func getObsoleteFiles(whitelist []string) []string {
	files := []string{}
	filepath.Walk(config.HSSHHostFolderPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == true {
			return nil
		}

		match := false
		for _, f := range whitelist {
			if f == path {
				match = true
			}
		}

		if match == true {
			return nil
		}

		files = append(files, path)

		return nil
	})

	return files
}

func syncWithProvider(providerConnection string) []string {
	projectID, remotePath, err := getProjectIDAndPath(providerConnection)

	provider, err := providers.New(providerConnection)
	if err != nil {
		messages.ProviderError(providerConnection, err)
		os.Exit(1)
	}

	files, err := provider.GetFiles(projectID, remotePath)
	if err != nil {
		messages.ProviderFetchError(providerConnection, err)
		os.Exit(1)
	}

	var wg = new(sync.WaitGroup)

	var filesCreated []string

	for _, file := range files {

		fileID := file.ID
		filePath := file.Path

		wg.Add(1)
		go func(path string, id string) {
			var hostPath string
			defer func() {
				if hostPath != "" {
					filesCreated = append(filesCreated, hostPath)
				}
				wg.Done()
			}()

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

			hostPath = host.GetPath()

		}(filePath, fileID)
	}

	wg.Wait()

	return filesCreated
}

// Sync ...
func Sync() {
	singleProvider := viper.GetString("provider")
	multiProvider := viper.GetStringSlice("providers")

	multiProvider = append(multiProvider, singleProvider)
	multiProvider = unique(multiProvider)

	var wg = new(sync.WaitGroup)

	var filesCreated []string

	for _, provider := range multiProvider {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if p == "" {
				return
			}

			filesCreated = append(filesCreated, syncWithProvider(p)...)
		}(provider)
	}

	wg.Wait()

	for _, file := range filesCreated {
		messages.SyncFileCreation(file)
	}

	obsoleteFiles := getObsoleteFiles(filesCreated)
	for _, file := range obsoleteFiles {
		messages.SyncFileDeletion(file)
		err := os.Remove(file)
		if err != nil {
			messages.CannotDeleteFile(err.Error(), file)
		}
	}
}
