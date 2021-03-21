package controllers

import (
	"hssh/config"
	"hssh/messages"
	"hssh/models"
	"os"
	"path/filepath"
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
	var filesToRead = []string{config.SSHConfigFilePath}

	go waitForParsedConnections(&connections, &channel)

	sshConfigInstance := models.NewSSHConfig(config.SSHConfigFilePath)
	filesToInclude := sshConfigInstance.GetIncludes()

	var folders = []string{
		config.HSSHHostFolderPath,
	}
	folders = unique(append(folders, filesToInclude...))

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

	filesToRead = unique(filesToRead)

	for _, file := range filesToRead {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			h := models.NewHost(f)
			h.ReadFile()
			h.Parse(&channel)
		}(file)
	}

	wg.Wait()

	time.Sleep(10 * time.Millisecond)

	if len(connections) == 0 {
		messages.NoConnections(connections)
		os.Exit(0)
	}

	// Sort alphabetically (case insensitive).
	sort.Slice(connections[:], func(i, j int) bool {
		return strings.ToLower(connections[i].Name) < strings.ToLower(connections[j].Name)
	})

	return connections
}
