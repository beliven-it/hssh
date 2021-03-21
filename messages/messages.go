package messages

import (
	"fmt"
	"hssh/models"
	"os"

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

// NoConfiguredYet ...
func NoConfiguredYet() {
	fmt.Println(Color("black", "\nHSSH configuration is not finished. Complete the step running the command:\n"))
	fmt.Println(Color("green", "\nhssh init"))
}

// MustBeConfigured ...
func MustBeConfigured() {
	fmt.Println(Color("yellow", "NOTE!"), Color("black", "The file must be configured before using the CLI"))
}

// ViperLoadError ...
func ViperLoadError(err error) {
	fmt.Println(Color("black", "Error reading the configuration file:"))
	fmt.Println(Color("red", err.Error()))
	os.Exit(1)
}

// PrintStep ...
func PrintStep(message string, err error) {
	status := Color("green", "OK")
	if err != nil {
		status = Color("red", "NOK")
	}

	stepStatus := Color("black", "[") + status + Color("black", "]")

	fmt.Println(stepStatus, Color("white", " "+message))
	if err != nil {
		fmt.Println(Color("red", err.Error()))
		os.Exit(1)
	}
}
