package controllers

import (
	"fmt"
	"hssh/config"
	"hssh/models"
	"io/ioutil"
	"os"
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
		fmt.Printf("File reading error: %v\n", err)
		os.Exit(1)
	}

	parseHostFile(string(data), channel)
}

func readHostFolder(path string, channel *chan models.Connection, wg *sync.WaitGroup) {
	defer wg.Done()

	var fs = new(sync.WaitGroup)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading files in folder: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		fs.Add(1)
		go readHostFile(path+"/"+file.Name(), channel, fs)
	}

	fs.Wait()
}

func translateToAbsolutePath(path string) string {
	if string(path[0]) == "/" {
		return path
	}

	return config.SSHFolderPath + "/" + path
}

func isFileIncluded(path string) bool {
	lastChar := string(path[len(path)-1])
	if lastChar == "*" || lastChar == "/" {
		return false
	}

	return true
}

func readConfig(channel *chan models.Connection) ([]string, []string) {
	includeRgx := regexp.MustCompile("(?m)^Include (.*)$")

	content, err := ioutil.ReadFile(config.SSHConfigFilePath)
	if err != nil {
		fmt.Printf("File reading error: %v\n", err)
		os.Exit(1)
	}

	body := string(content)
	includes := includeRgx.FindAllStringSubmatch(body, -1)

	includesFoldersPath := []string{}
	includesFilesPath := []string{}
	for _, include := range includes {
		pathTranslated := translateToAbsolutePath(include[1])
		if isFileIncluded(pathTranslated) {
			includesFilesPath = append(includesFilesPath, pathTranslated)
			continue
		}

		pathTranslated = strings.Replace(pathTranslated, "*", "", -1)
		includesFoldersPath = append(includesFoldersPath, pathTranslated)
	}

	parseHostFile(body, channel)

	return includesFoldersPath, includesFilesPath
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

	var folders = []string{
		config.HSSHHostFolderPath + "/",
	}

	go waitForParsedConnections(&connections, &channel)

	foldersToInclude, filesToRead := readConfig(&channel)

	folders = unique(append(folders, foldersToInclude...))

	for _, file := range filesToRead {
		wg.Add(1)
		go readHostFile(file, &channel, wg)
	}

	for _, folder := range folders {
		wg.Add(1)
		go readHostFolder(folder, &channel, wg)
	}

	wg.Wait()

	time.Sleep(10 * time.Millisecond)

	// Sort alphabetically (case insensitive).
	sort.Slice(connections[:], func(i, j int) bool {
		return strings.ToLower(connections[i].Name) < strings.ToLower(connections[j].Name)
	})

	return connections
}
