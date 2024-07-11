package services

import (
	"hssh/models"
	"os"
)

// Find ...
func Find(host string) models.Connection {
	connections := List()

	if host == "" {
		commandVerbToExec := serializeConnections(&connections)

		// Choose connection
		host = fzf(commandVerbToExec)
		if host == "" {
			os.Exit(0)
		}
	}

	return searchConnectionByPattern(host, &connections)
}
