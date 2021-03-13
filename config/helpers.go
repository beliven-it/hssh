package config

import (
	"os"
)

func getHomePath() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}

	return homePath
}
