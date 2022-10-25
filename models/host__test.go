package models

import (
	"testing"
)

func TestHostParse(t *testing.T) {
	hostRaw := `Host test
    Hostname 192.168.1.1
    User johndoe
    Port 1234`

	i := NewHost("")

	connections := i.ParseRow(hostRaw)

	connection := connections[0]

	if connection.Name != "test" {
		t.Errorf("The connection must have the name test and not %s", connection.Name)
	}
	if connection.Hostname != "192.168.1.1" {
		t.Errorf("The connection must have the hostname 192.168.1.1 and not %s", connection.Hostname)
	}
	if connection.User != "johndoe" {
		t.Errorf("The connection must have the user johndoe and not %s", connection.User)
	}
	if connection.Port != "1234" {
		t.Errorf("The connection must have the port 1234 and not %s", connection.Port)
	}
}

func TestHostParseUsingHostName(t *testing.T) {
	hostRaw := `Host test
    HostName 192.168.1.1
    User johndoe
    Port 1234`

	i := NewHost("")
	connections := i.ParseRow(hostRaw)

	connection := connections[0]

	if connection.Hostname != "192.168.1.1" {
		t.Errorf("The connection must have the hostname 192.168.1.1 and not %s", connection.Hostname)
	}
}

func TestHOstParseUsingHostNameThatContainsPort(t *testing.T) {
	hostRaw := `Host testporton
    HostName 192.168.1.1
    User johndoe
    Port 1234`

	i := NewHost("")
	connections := i.ParseRow(hostRaw)

	connection := connections[0]

	if connection.Name != "testporton" {
		t.Errorf("The connection must have the host set to testporton and not %s", connection.Hostname)
	}
}

func TestHostUsingMultipleAliases(t *testing.T) {
	hostRaw := `Host pearl ruby zapphire "stone"
    HostName 192.168.1.1
    User johndoe
    Port 1234`

	i := NewHost("")
	connections := i.ParseRow(hostRaw)

	if len(connections) != 4 {
		t.Errorf("The connections count are 4 and not %d", len(connections))
	}

	if connections[0].Name != "pearl" {
		t.Errorf("The first connections name must named pearl and not %s", connections[0].Name)
	}

	if connections[3].Name != "\"stone\"" {
		t.Errorf("The last connections name must named \"stone\" and not %s", connections[3].Name)
	}
}
