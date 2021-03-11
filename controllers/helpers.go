package controllers

import (
	"bytes"
	"os"
	"os/exec"
)

func fzf(context string) string {
	cmdOutput := &bytes.Buffer{}
	c := exec.Command("bash", "-c", "echo -e '"+context+"' | fzf")
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
