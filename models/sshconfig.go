package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type sshconfig struct {
	filesToInclude []string
	content        string
	path           string
}

// ISSHConfig ...
type ISSHConfig interface {
	GetIncludes() []string
	GetContent() string
	GetPath() string
	SetContent(string)
}

func (s *sshconfig) translatePath(path string) string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	firstChar := string(path[0])

	if firstChar == "/" {
		return path
	}

	if firstChar == "~" {
		return homePath + string(path[1:])
	}

	return homePath + "/.ssh/" + path
}

func (s *sshconfig) takeIncludes() []string {
	rgx := regexp.MustCompile("(?m)^(\\s+|)Include (.*)$")
	includes := rgx.FindAllStringSubmatch(s.content, -1)

	includesList := []string{}
	for _, include := range includes {
		includesList = append(includesList, s.translatePath(include[2]))
	}

	return includesList
}

func (s *sshconfig) readFile() {
	content, err := ioutil.ReadFile(s.path)
	if err != nil {
		fmt.Printf("File reading error: %v\n", err)
		// TODO: Remove this exit and handle it in final side
		os.Exit(1)
	}

	s.SetContent(string(content))
	s.filesToInclude = s.takeIncludes()
}

func (s *sshconfig) GetContent() string {
	return s.content
}

func (s *sshconfig) GetIncludes() []string {
	return s.filesToInclude
}

func (s *sshconfig) GetPath() string {
	return s.path
}

func (s *sshconfig) SetContent(content string) {
	s.content = content
}

// NewSSHConfig ...
func NewSSHConfig(path string) ISSHConfig {
	s := sshconfig{
		path: path,
	}

	s.readFile()

	return &s
}
