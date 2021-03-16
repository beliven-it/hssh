package controllers

import (
	"bytes"
	"hssh/models"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func serializeConnections(connections *[]models.Connection) string {
	// Format in one string
	listOfConnectionsNames := []string{}
	for _, connection := range *connections {
		listOfConnectionsNames = append(listOfConnectionsNames, connection.Name+" -> "+connection.Hostname)
	}

	return strings.Join(listOfConnectionsNames, "\n")
}

func fzf(context string) string {
	cmdOutput := &bytes.Buffer{}
	c := exec.Command("bash", "-c", "echo -e '"+context+"' | fzf")
	c.Stdout = cmdOutput
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	err := c.Run()
	if err != nil && err.Error() != "exit status 130" {
		os.Exit(1)
	}

	return string(cmdOutput.Bytes())
}

func ssh(connectionName string) {
	c := exec.Command("bash", "-c", "ssh "+connectionName)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	err := c.Run()
	if err != nil {
		os.Exit(1)
	}
}

func fromFzfSelectionToConnection(selection string, connections *[]models.Connection) models.Connection {
	clearRgx := regexp.MustCompile("(\\n|\\r)")
	takeNameRgx := regexp.MustCompile(" ->.*")
	selection = clearRgx.ReplaceAllString(selection, "")
	selection = takeNameRgx.ReplaceAllString(selection, "")
	for _, connection := range *connections {
		if connection.Name == selection {
			return connection
		}
	}

	return models.Connection{}
}
