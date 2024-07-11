package services

import (
	"os"
)

// Connect ...
func Connect(host string) {
	if host == "" {
		// Select the connection using FZF
		connections := List()

		commandVerbToExec := serializeConnections(&connections)

		// Choose connection
		connectionString := fzf(commandVerbToExec)
		if connectionString == "" {
			os.Exit(0)
		}

		connection := searchConnectionByPattern(connectionString, &connections)

		host = connection.Name
	}

	// connect via ssh
	ssh(host)
}
