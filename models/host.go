package models

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// Host ...
type host struct {
	path    string
	content string
}

// IHost ...
type IHost interface {
	ReadFile()
	Parse(string) Connection
	List(*chan Connection)
	GetPath() string
	GetContent() string
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

func (h *host) Parse(hostRaw string) Connection {
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
		} else {
			connection.Name = strings.ReplaceAll(attribute, "Host ", "")
		}
	}

	return connection
}

func (h *host) List(channel *chan Connection) {
	content := strings.TrimSpace(h.content)

	// Remove comments
	content = regexp.MustCompile("(?m)^#.*").ReplaceAllString(content, "")

	// Remove empty lines
	content = regexp.MustCompile("[\t\r\n]+").ReplaceAllString(content, "\n")

	// Split content into hosts
	hosts := strings.Split(content, "Host ")

	for x, host := range hosts {
		if x == 0 {
			continue
		}

		connection := h.Parse(host)

		if channel != nil {
			go func() {
				*channel <- connection
			}()
		}
	}

}

func (h *host) GetPath() string {
	return h.path
}

func (h *host) GetContent() string {
	return h.content
}

// NewHost ...
func NewHost(path string) IHost {
	return &host{
		path: path,
	}
}
