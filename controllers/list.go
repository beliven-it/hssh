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

func isBlackListed(hostname string, blacklist []string) bool {
	for _, b := range blacklist {
		if b == hostname {
			return true
		}
	}

	return false
}

func waitForParsedConnections(connections *[]models.Connection, channel *chan models.Connection, wg *sync.WaitGroup) {
	for connection := range *channel {
		if isBlackListed(connection.Name, []string{"*"}) == false {
			*connections = append(*connections, connection)
		}
		wg.Done()
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

	conns := [][]models.Connection{}
	chans := []chan models.Connection{}

	for i, file := range filesToRead {
		conns = append(conns, []models.Connection{})
		chans = append(chans, make(chan models.Connection))

		go func(f string, ch *chan models.Connection, index int) {
			go waitForParsedConnections(&conns[index], ch, wg)
			h := models.NewHost(f)
			h.ReadFile()
			h.Parse()
			wg.Add(h.GetConnectionsCount())

			h.ProvideViaChannel(ch)
		}(file, &chans[i], i)
	}

	time.Sleep(10 * time.Millisecond)
	wg.Wait()

	for i := 0; i < len(conns); i++ {
		connections = append(connections, conns[i]...)
		close(chans[i])
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
