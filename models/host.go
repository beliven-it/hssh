package models

import (
	"os"
	"regexp"
	"strings"
)

// Host ...
type host struct {
	path        string
	content     string
	connections []Connection
}

// IHost ...
type IHost interface {
	ReadFile()
	ParseRow(string) []Connection
	Parse() []Connection
	GetPath() string
	GetContent() string
	GetConnectionsCount() int
	GetConnections() []Connection
	ProvideViaChannel(*chan Connection)
	Create([]byte) error
}

func (h *host) ReadFile() {
	contentBytes, err := os.ReadFile(h.path)
	if err != nil {
		return
	}

	h.content = string(contentBytes)
}

func (h *host) Create(content []byte) error {
	filePathSplitted := strings.Split(h.path, "/")
	folderPath := filePathSplitted[0 : len(filePathSplitted)-1]
	os.MkdirAll(strings.Join(folderPath, "/"), os.ModePerm)

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

func (h *host) getAttributeFromRow(attribute string, row string) string {
	rgx := regexp.MustCompile(attribute + "(|\\s)")
	return strings.Trim(rgx.ReplaceAllString(row, ""), " ")
}

func (h *host) getAliases(host string) []string {
	var aliases = []string{}
	rgxDoubleQuotes := regexp.MustCompile(`(".*?")`)

	matches := rgxDoubleQuotes.FindAllStringSubmatch(host, -1)
	if len(matches) == 0 {
		aliases = append(aliases, host)
		return aliases
	}

	host = rgxDoubleQuotes.ReplaceAllString(host, "")
	aliases = append(aliases, strings.Split(host, " ")...)

	for _, group := range matches {
		aliases = append(aliases, group[1:]...)
	}

	filteredAliases := []string{}
	for _, alias := range aliases {
		alias = strings.Trim(alias, " ")
		if alias == "" {
			continue
		}

		if alias == "*" {
			continue
		}

		filteredAliases = append(filteredAliases, alias)
	}

	return filteredAliases
}

func (h *host) ParseRow(hostRaw string) []Connection {
	connection := Connection{}

	for _, attribute := range strings.Split(hostRaw, "\n") {
		attribute = strings.Trim(attribute, " ")
		partials := strings.Split(attribute, " ")

		if partials[0] == "" {
			continue
		}

		attribute = strings.ToLower(partials[0]) + " " + partials[1]

		if strings.HasPrefix(attribute, "hostname ") {
			connection.Hostname = h.getAttributeFromRow("hostname", attribute)
		}

		if strings.HasPrefix(attribute, "user ") {
			connection.User = h.getAttributeFromRow("user", attribute)
		}

		if strings.HasPrefix(attribute, "port ") {
			connection.Port = h.getAttributeFromRow("port", attribute)
		}

		if strings.HasPrefix(attribute, "identityfile ") {
			connection.IdentityFile = h.getAttributeFromRow("identityfile", attribute)
		}

		if strings.HasPrefix(attribute, "host ") {
			connection.Name = h.getAttributeFromRow("host", attribute)
		}
	}

	aliases := h.getAliases(connection.Name)

	var connections = []Connection{}

	if !connection.IsWellConfigured() {
		return h.connections
	}

	for _, alias := range aliases {
		connection.Name = alias

		if connection.IsAllowed() {
			connections = append(connections, connection)
		}
	}

	return connections
}

func (h *host) ProvideViaChannel(channel *chan Connection) {
	if channel == nil {
		return
	}

	for _, connection := range h.connections {
		go func(c Connection) {
			*channel <- c
		}(connection)
	}
}

func (h *host) Parse() []Connection {
	content := strings.TrimSpace(h.content)

	// Remove comments
	content = regexp.MustCompile(`(?m)^(|\s+)#.*`).ReplaceAllString(content, "")

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

		h.connections = append(h.connections, h.ParseRow(host)...)
	}

	return h.connections
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

func (h *host) GetConnectionsCount() int {
	return len(h.connections)
}

// NewHost ...
func NewHost(path string) IHost {
	return &host{
		path: path,
	}
}
