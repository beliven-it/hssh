package controllers

import (
	"hssh/models"
	"os"
)

// Find ...
func Find() models.Connection {
	connections := List()

	commandVerbToExec := serializeConnections(&connections)

	// Choose connection
	connectionString := fzf(commandVerbToExec)
	if connectionString == "" {
		os.Exit(0)
	}

	return fromFzfSelectionToConnection(connectionString, &connections)
}
