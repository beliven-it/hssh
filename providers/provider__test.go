package providers

import (
	"testing"
)

// TestGetDriverFromConnectionString ...
func TestGetDriverFromConnectionString(t *testing.T) {
	driver, e := getDriverFromConnectionString("github://12345678://test")
	if e != nil {
		t.Errorf("Should not return any error")
	}

	if driver != "github" {
		t.Errorf("Should return the first part of the string instead of %s", driver)
	}

	driver, e = getDriverFromConnectionString("gitlab://:/test")
	if e != nil {
		t.Errorf("Should not return any error")
	}

	if driver != "gitlab" {
		t.Errorf("Should return the first part of the string instead of %s", driver)
	}
}

func withFakeProvider(driver string, connectionString string) provider {
	p := provider{
		connectionString: connectionString,
	}
	p.ParseConnection(driver)
	return p
}

// TestParseConnection ...
func TestParseConnection(t *testing.T) {
	r := withFakeProvider("github", "github://123456:/CasvalDOT/hssh@providers")
	pt := r.GetPrivateToken()

	if pt != "123456" {
		t.Errorf("Should return the provided token instead of %s", pt)
	}

	r = withFakeProvider("github", "github://:/CasvalDOT/hssh@providers")
	pt = r.GetPrivateToken()

	if pt != "" {
		t.Errorf("Should return an empty private token")
	}
}
