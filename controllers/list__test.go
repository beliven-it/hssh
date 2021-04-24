package controllers

import (
	"testing"
)

func TestBlacklist(t *testing.T) {
	isBlacklist := isBlackListed("anna", []string{"*"})
	if isBlacklist == true {
		t.Errorf("The hostname cannot be blacklisted")
	}
	isBlacklist = isBlackListed("*", []string{"*"})
	if isBlacklist == false {
		t.Errorf("The hostname must be blacklisted")
	}
	isBlacklist = isBlackListed("", []string{"*", ""})
	if isBlacklist == false {
		t.Errorf("The hostname must be blacklisted")
	}
	isBlacklist = isBlackListed(" ", []string{"*", ""})
	if isBlacklist == false {
		t.Errorf("The hostname must be blacklisted")
	}
}
