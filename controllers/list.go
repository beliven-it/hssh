package controllers

import (
	"hssh/config"
	"hssh/messages"
	"hssh/models"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func unique(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}

	for e := range arr {
		if !occured[arr[e]] {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}

	return result
}

// List the connections available
func List() []models.Connection {
	var connections []models.Connection
	var filesToRead = []string{config.SSHConfigFilePath}

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
		h := models.NewHost(file)
		h.ReadFile()
		connections = append(connections, h.Parse()...)
	}

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
