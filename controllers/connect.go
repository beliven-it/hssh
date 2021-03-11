package controllers

import (
	"fmt"
	"os"
	"strings"
)

// Connect ...
func Connect() {
	// Select the connection using FZF
	connections := List()

	// Format in one string
	listOfConnectionsNames := []string{}
	for _, connection := range connections {
		listOfConnectionsNames = append(listOfConnectionsNames, connection.Name)
	}

	commandVerbToExec := strings.Join(listOfConnectionsNames, "\n")

	// Choose connection
	connectionName := fzf(commandVerbToExec)
	if connectionName == "" {
		fmt.Println("Selection is empty. The request is rejected")
		os.Exit(1)
	}

	// connect via ssh
	ssh(connectionName)

}
