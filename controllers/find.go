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
	connectionName := fzf(commandVerbToExec)
	if connectionName == "" {
		os.Exit(0)
	}

	return fromFzfSelectionToConnection(connectionName, &connections)
}
