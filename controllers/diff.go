package controllers

import (
	"hssh/messages"
	"hssh/providers"
	"os"

	"github.com/spf13/viper"
)

func checkProvider(providerConnection string) {
	projectID, remotePath, err := getProjectIDAndPath(providerConnection)

	provider, err := providers.New(providerConnection)
	if err != nil {
		messages.ProviderError(providerConnection, err)
		os.Exit(1)
	}
}

// Diff check for each provider if any changes
// remotely occured
func Diff() {
	// For each provider in config
	// file iterate and if something changes
	// at repository (or not repository), the user must be notifier
	providers := viper.GetStringSlice("providers")
	provider := viper.GetString("povider")

	// Concatenate the single provider option
	// into the multiprovider array
	providers = append(providers, provider)
	providers = unique(providers)

	for _, p := range providers {
		checkProvider(p)
	}
}
