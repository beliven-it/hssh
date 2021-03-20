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
	connection := i.Parse(hostRaw)

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
