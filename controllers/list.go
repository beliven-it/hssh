package controllers

import (
	"fmt"
	"hssh/config"
	"hssh/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

func waitForParsedConnections(connections *[]models.Connection, channel *chan models.Connection) {
	for connection := range *channel {
		*connections = append(*connections, connection)
	}
}

func readHostFile(path string, channel *chan models.Connection, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	parseHostFile(string(data), channel)
}

func translateToAbsolutePath(path string) string {
	firstChar := string(path[0])
	if firstChar == "/" {
		return path
	} else if firstChar == "~" {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return path
		}

		path = string(path[1:])
		return homePath + path
	} else {
		path = config.SSHFolderPath + "/" + path
	}

	stat, err := os.Stat(path)
	if err != nil {
		return path
	}

	if stat.IsDir() == true {
		path = path + "/*"
	}

	return path
}

func isFileIncluded(path string) bool {
	lastChar := string(path[len(path)-1])
	if lastChar == "*" || lastChar == "/" {
		return false
	}

	return true
}

func readConfig(channel *chan models.Connection) []string {
	includeRgx := regexp.MustCompile("(?m)^Include (.*)$")

	content, err := ioutil.ReadFile(config.SSHConfigFilePath)
	if err != nil {
		fmt.Printf("File reading error: %v\n", err)
		os.Exit(1)
	}

	body := string(content)
	includes := includeRgx.FindAllStringSubmatch(body, -1)

	includesPath := []string{}
	for _, include := range includes {
		pathTranslated := translateToAbsolutePath(include[1])

		includesPath = append(includesPath, pathTranslated)
	}

	parseHostFile(body, channel)

	return includesPath
}

func unique(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}

	for e := range arr {
		if occured[arr[e]] != true {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}

	return result
}

// List the connections available
func List() []models.Connection {
	var wg = new(sync.WaitGroup)
	var channel = make(chan models.Connection)
	var connections []models.Connection
	var filesToRead = []string{}

	var folders = []string{
		config.HSSHHostFolderPath,
	}

	go waitForParsedConnections(&connections, &channel)

	foldersToInclude := readConfig(&channel)

	folders = unique(append(folders, foldersToInclude...))

	// Take all files from folder
	for _, folder := range folders {
		files, ok := filepath.Glob(folder)
		if ok != nil {
			continue
		}

		if len(files) == 0 {
			filesToRead = append(filesToRead, folder)
		}

		filesToRead = append(filesToRead, files...)
	}

	for _, file := range filesToRead {
		wg.Add(1)
		go readHostFile(file, &channel, wg)
	}

	wg.Wait()

	time.Sleep(10 * time.Millisecond)

	// Sort alphabetically (case insensitive).
	sort.Slice(connections[:], func(i, j int) bool {
		return strings.ToLower(connections[i].Name) < strings.ToLower(connections[j].Name)
	})

	return connections
}
