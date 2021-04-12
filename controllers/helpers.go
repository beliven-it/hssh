package controllers

import (
	"bytes"
	"fmt"
	"hssh/models"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// PrintConnectionDetails ...
func PrintConnectionDetails(connection *models.Connection) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgHiMagenta).SprintFunc()

	fmt.Printf(
		"\nName: %s\nHostname: %s\nUser: %s\nPort: %s\nIdentity: %s\n",
		green(connection.Name), magenta(connection.Hostname), blue(connection.User), red(connection.Port), yellow(connection.IdentityFile),
	)
}

// PrintConnection ...
func PrintConnection(connection *models.Connection, withColors bool) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	if withColors {
		fmt.Printf("%s -> %s@%s:%s %s\n", green(connection.Name), blue(connection.User), connection.Hostname, red(connection.Port), yellow(connection.IdentityFile))
	} else {
		fmt.Printf("%s -> %s@%s:%s %s\n", connection.Name, connection.User, connection.Hostname, connection.Port, connection.IdentityFile)
	}
}

func serializeConnections(connections *[]models.Connection) string {
	listOfConnectionsNames := []string{}
	for _, connection := range *connections {
		listOfConnectionsNames = append(listOfConnectionsNames, connection.Name+" -> "+connection.Hostname)
	}

	return strings.Join(listOfConnectionsNames, "\n")
}

func fzf(context string) string {
	cmdOptions := viper.GetString("fzf_options")
	cmdOutput := &bytes.Buffer{}
	c := exec.Command("bash", "-c", "echo -e '"+context+"' | fzf "+cmdOptions)
	c.Stdout = cmdOutput
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	err := c.Run()
	if err != nil && err.Error() != "exit status 130" {
		os.Exit(1)
	}

	return string(cmdOutput.Bytes())
}

func ssh(connectionName string) {
	c := exec.Command("bash", "-c", "ssh "+connectionName)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	err := c.Run()
	if err != nil {
		os.Exit(1)
	}
}

func searchConnectionByPattern(pattern string, connections *[]models.Connection) models.Connection {
	clearRgx := regexp.MustCompile("(\\n|\\r)")
	takeNameRgx := regexp.MustCompile(" ->.*")
	pattern = clearRgx.ReplaceAllString(pattern, "")
	pattern = takeNameRgx.ReplaceAllString(pattern, "")
	for _, connection := range *connections {
		if connection.Name == pattern {
			return connection
		}
	}

	return models.Connection{}
}
