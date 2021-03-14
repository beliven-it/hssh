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
)

// List the connections available
func List() []models.Connection {
	files, err := ioutil.ReadDir(config.HSSHHostFolderPath)
	if err != nil {
		fmt.Printf("Error reading files in folder: %v\n", err)
		os.Exit(1)
	}

	content := ""
	for _, file := range files {
		data, err := ioutil.ReadFile(config.HSSHHostFolderPath + "/" + file.Name())
		if err != nil {
			fmt.Printf("File reading error: %v\n", err)
			os.Exit(1)
		}

		// Convert byte to string and add to content.
		content += string(data)
	}

	// Remove comments from content.
	content = regexp.MustCompile("(?m)^#.*").ReplaceAllString(content, "")

	// Remove empty lines from content.
	content = regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(content), "\n")

	// Split content into hosts.
	hosts := strings.Split(content, "Host ")

	// Map hosts into array of Connection struct.
	var connections []models.Connection
	for indexHost, host := range hosts {
		if indexHost == 0 {
			continue
		}

		host = strings.ReplaceAll(host, " ", "")

		var temp = models.Connection{}
		for indexParam, param := range strings.Split(host, "\n") {

			if indexParam == 0 {
				temp.Name = param
			} else {
				if strings.Contains(param, "Hostname") {
					temp.Hostname = strings.ReplaceAll(param, "Hostname", "")
				}

				if strings.Contains(param, "User") {
					temp.User = strings.ReplaceAll(param, "User", "")
				}

				if strings.Contains(param, "Port") {
					temp.Port = strings.ReplaceAll(param, "Port", "")
				}

				if strings.Contains(param, "IdentityFile") {
					temp.IdentityFile = strings.ReplaceAll(param, "IdentityFile", "")
				}
			}
		}

		connections = append(connections, temp)
	}

	// Sort alphabetically (case insensitive).
	sort.Slice(connections[:], func(i, j int) bool {
		return strings.ToLower(connections[i].Name) < strings.ToLower(connections[j].Name)
	})

	return connections
}
