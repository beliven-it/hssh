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
	connectionName := fzf(commandVerbToExec)
	if connectionName == "" {
		fmt.Println("Selection is empty. The request is rejected")
		os.Exit(1)
	}

	// connect via ssh
	ssh(connectionName)
}
