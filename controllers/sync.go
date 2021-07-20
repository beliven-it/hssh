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

	wg := new(sync.WaitGroup)

	wg.Add(len(files))
	var filesCreated []string
	channel := make(chan string, len(files))

	for _, file := range files {
		go func(path string, id string) {
			fileContent, err := provider.GetFile(projectID, id)
			var hostPath string
			defer func(hp string) {
				channel <- hostPath
			}(hostPath)

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
			return
		}(file.Path, file.ID)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for c := range channel {
		if c != "" {
			filesCreated = append(filesCreated, c)
		}
		wg.Done()
	}

	return filesCreated
}

// Sync ...
func Sync() {
	singleProvider := viper.GetString("provider")
	multiProvider := viper.GetStringSlice("providers")

	multiProvider = append(multiProvider, singleProvider)
	multiProvider = unique(multiProvider)

	wg := new(sync.WaitGroup)

	var filesCreated []string
	channel := make(chan []string, len(multiProvider))

	wg.Add(len(multiProvider))

	for _, provider := range multiProvider {
		go func(p string) {
			if p == "" {
				channel <- []string{}
			} else {
				channel <- syncWithProvider(p)
			}
			return
		}(provider)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for c := range channel {
		filesCreated = append(filesCreated, c...)
		wg.Done()
	}

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
