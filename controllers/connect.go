package controllers

import (
	"fmt"
	"os"
)

// Connect ...
func Connect() {
	// Select the connection using FZF
	connections := List()

	commandVerbToExec := serializeConnections(&connections)

	// Choose connection
	connectionString := fzf(commandVerbToExec)
	if connectionString == "" {
		fmt.Println("Selection is empty. The request is rejected")
		os.Exit(1)
	}

	connection := fromFzfSelectionToConnection(connectionString, &connections)

	// connect via ssh
	ssh(connection.Name)
}
