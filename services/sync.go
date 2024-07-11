package services

import (
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

type Config struct {
	Providers []providers.ProviderConnection `mapstructure:"providers"`
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
		if info.IsDir() {
			return nil
		}

		match := false
		for _, f := range whitelist {
			if f == path {
				match = true
			}
		}

		if match {
			return nil
		}

		files = append(files, path)

		return nil
	})

	return files
}

func syncWithProvider(providerConnection providers.ProviderConnection) []string {
	provider, err := providers.New(providerConnection)
	if err != nil {
		os.Exit(1)
	}

	files, err := provider.GetFiles(providerConnection.EntityID, providerConnection.Subpath)
	if err != nil {
		messages.ProviderFetchError(err)
		os.Exit(1)
	}

	wg := new(sync.WaitGroup)

	wg.Add(len(files))
	var filesCreated []string
	channel := make(chan string, len(files))

	for _, file := range files {
		go func(path string, id string) {
			fileContent, err := provider.GetFile(providerConnection.EntityID, id)
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

func readProviderConnections() []providers.ProviderConnection {
	var providersConfig Config
	list := []providers.ProviderConnection{}

	// Single provider of the first version
	singleProvider := viper.GetString("provider")

	// Multiple provider of the early version
	// if the user select the new version of multi structured providers
	// the result is empty slice
	multiProvider := viper.GetStringSlice("providers")
	multiProvider = append(multiProvider, singleProvider)
	multiProvider = unique(multiProvider)

	// Structured multi provider for the new version
	viper.Unmarshal(&providersConfig)

	// Convert string connections into structured version
	for _, p := range multiProvider {
		pconn := providers.ProviderConnection{}
		pconn.FromString(p)

		list = append(list, pconn)
	}

	return append(list, providersConfig.Providers...)
}

// Sync ...
func Sync() {
	multiProvider := readProviderConnections()

	wg := new(sync.WaitGroup)

	var filesCreated []string
	channel := make(chan []string, len(multiProvider))

	wg.Add(len(multiProvider))

	for _, provider := range multiProvider {
		go func(p providers.ProviderConnection) {
			if p.Type == "" {
				channel <- []string{}
			} else {
				channel <- syncWithProvider(p)
			}
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
