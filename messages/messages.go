package messages

import (
	"fmt"
	"hssh/models"
	"os"
	"regexp"

	"github.com/fatih/color"
)

// Color ...
func Color(c string, content string) string {
	switch c {
	case "red":
		return color.New(color.FgRed).SprintFunc()(content)
	case "green":
		return color.New(color.FgGreen).SprintFunc()(content)
	case "yellow":
		return color.New(color.FgYellow).SprintFunc()(content)
	case "magenta":
		return color.New(color.FgMagenta).SprintFunc()(content)
	case "black":
		return color.New(color.FgHiBlack).SprintFunc()(content)
	case "blue":
		return color.New(color.FgBlue).SprintFunc()(content)
	}

	return content
}

// NoConnections ...
func NoConnections(connections []models.Connection) {
	if len(connections) == 0 {
		fmt.Println(Color("black", "\nThere aren't host connections."))
		fmt.Println(Color("black", "Try to run:\n\n"), Color("green", "hssh sync"))
		fmt.Println(Color("black", "\nIf doesn't work check the ") + Color("yellow", "~/.ssh/config") + Color("black", " and make sure there is at least one host connection configured"))
	}
}

// NoConnection ...
func NoConnection() {
	fmt.Println(Color("black", "\nCannot find the details of the host!"))
	fmt.Println(Color("black", "Retry using "+Color("blue", "fuzzysearch")) + "\n")
	fmt.Println(Color("green", "hssh find"))
}

// ConfigNotEditedYet ...
func ConfigNotEditedYet() {
	fmt.Println(Color("black", "\nHSSH configuration is not yet finished. Take time to configure the provider option\n"))
	fmt.Println(Color("green", "~/.config/hssh/config.yml"))
	fmt.Println(Color("black", "\nIf file doesn't exist run:\n"))
	fmt.Println(Color("green", "hssh init"))
}

// NoConfiguredYet ...
func NoConfiguredYet() {
	fmt.Println(Color("black", "\nHSSH configuration is not finished. Complete the step running the command:\n"))
	fmt.Println(Color("green", "hssh init"))
}

// MustBeConfigured ...
func MustBeConfigured() {
	fmt.Println(Color("yellow", "NOTE!"), Color("black", "The file must be configured before using the CLI"))
}

// ViperLoadError ...
func ViperLoadError(err error) {
	fmt.Println(Color("black", "Error reading the configuration file:"))
	fmt.Println(Color("red", err.Error()))
	fmt.Println(Color("black", "If config.yml is missing run:\n"))
	fmt.Println(Color("green", "hssh init\n"))
	os.Exit(1)
}

// ProviderError ...
func ProviderError(connectionString string, err error) {
	fmt.Println(Color("red", err.Error()))
	fmt.Println(Color("black", "Checkout the config file for errors."))
	fmt.Println(Color("black", "Probabily the connection string is malformed"))
	fmt.Println(Color("blue", "\n"+connectionString+"\n"))
}

//ProviderFetchError ...
func ProviderFetchError(connectionString string, err error) {
	fmt.Println(Color("black", "An error occured during files fetch"))
	ProviderError(connectionString, err)
}

// SyncFileCreation ...
func SyncFileCreation(path string) {
	rgx := regexp.MustCompile("(.*)\\.(.*)")
	pathMatchs := rgx.FindAllStringSubmatch(path, 1)
	body := pathMatchs[0][1]
	extension := pathMatchs[0][2]

	fmt.Println(Color("black", "[") + Color("green", "CREATED") + Color("black", "]") + " " + body + "." + Color("yellow", extension))
}

// SyncFileDeletion ...
func SyncFileDeletion(path string) {
	fmt.Println(Color("black", "[") + Color("red", "DELETED") + Color("black", "]") + " " + path)
}

// CannotDeleteFile ...
func CannotDeleteFile(message, file string) {
	fmt.Println(Color("black", "Cannot delete file"), Color("blue", file))
	fmt.Println(Color("red", message))
}

// Print ...
func Print(message string) {
	fmt.Println(Color("blue", message))
}

// PrintStep ...
func PrintStep(message string, err error) {
	status := Color("green", "OK")
	if err != nil {
		status = Color("red", "NOK")
	}

	stepStatus := Color("black", "[") + status + Color("black", "]")

	fmt.Println(stepStatus, Color("white", message))
	if err != nil {
		fmt.Println(Color("red", err.Error()))
		os.Exit(1)
	}
}
