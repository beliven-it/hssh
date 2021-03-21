package models

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// Host ...
type host struct {
	path        string
	content     string
	connections []Connection
	channel     chan Connection
}

// IHost ...
type IHost interface {
	ReadFile()
	ParseRow(string) Connection
	Parse(*chan Connection)
	GetPath() string
	GetContent() string
	GetConnections() []Connection
	Create([]byte) error
}

func (h *host) ReadFile() {
	contentBytes, err := ioutil.ReadFile(h.path)
	if err != nil {
		return
	}

	h.content = string(contentBytes)
}

func (h *host) Create(content []byte) error {
	file, err := os.Create(h.path)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}

func (h *host) ParseRow(hostRaw string) Connection {
	connection := Connection{}
	for _, attribute := range strings.Split(hostRaw, "\n") {
		attribute = strings.Trim(attribute, " ")

		if attribute == "" {
			continue
		}

		if strings.Contains(attribute, "Hostname") {
			connection.Hostname = strings.ReplaceAll(attribute, "Hostname ", "")
		} else if strings.Contains(attribute, "User") {
			connection.User = strings.ReplaceAll(attribute, "User ", "")
		} else if strings.Contains(attribute, "Port") {
			connection.Port = strings.ReplaceAll(attribute, "Port ", "")
		} else if strings.Contains(attribute, "IdentityFile") {
			connection.IdentityFile = strings.ReplaceAll(attribute, "IdentityFile ", "")
		} else if strings.Contains(attribute, "Host ") {
			connection.Name = strings.ReplaceAll(attribute, "Host ", "")
		}
	}

	return connection
}

func (h *host) provideViaChannel(channel *chan Connection) {
	if channel == nil {
		return
	}

	for _, connection := range h.connections {
		go func(c Connection) {
			*channel <- c
		}(connection)

	}
}

func (h *host) Parse(channel *chan Connection) {
	content := strings.TrimSpace(h.content)

	// Remove comments
	content = regexp.MustCompile("(?m)^#.*").ReplaceAllString(content, "")

	// Remove empty lines
	content = regexp.MustCompile("[\t\r\n]+").ReplaceAllString(content, "\n")

	// Apply a marker for splitting logic
	content = regexp.MustCompile("Host ").ReplaceAllString(content, "!!Host ")

	// Split content into hosts
	hosts := strings.Split(content, "!!")

	for x, host := range hosts {
		if x == 0 {
			continue
		}
		h.connections = append(h.connections, h.ParseRow(host))
	}

	go h.provideViaChannel(channel)
}

func (h *host) GetPath() string {
	return h.path
}

func (h *host) GetContent() string {
	return h.content
}

func (h *host) GetConnections() []Connection {
	return h.connections
}

// NewHost ...
func NewHost(path string) IHost {
	return &host{
		path: path,
	}
}
