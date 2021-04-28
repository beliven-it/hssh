package models

import (
	"testing"
)

func TestConnectionIsAllowed(t *testing.T) {
	connection := Connection{
		Name: "*",
	}
	isAllowed := connection.IsAllowed()
	if isAllowed {
		t.Errorf("The connection %s is not allowed", connection.Name)
	}

	connection = Connection{
		Name: "Heply",
	}
	isAllowed = connection.IsAllowed()
	if !isAllowed {
		t.Errorf("The connection %s is allowed", connection.Name)
	}
}
