package models

// Connection struct.
type Connection struct {
	Name         string
	Hostname     string
	User         string
	Port         string
	IdentityFile string
}

// IsWellConfigured ...
// Check if the connection is well configured
func (c *Connection) IsWellConfigured() bool {
	if c.Name == "" || c.Hostname == "" {
		return false
	}

	return true
}

// IsAllowed ...
func (c *Connection) IsAllowed() bool {
	return c.Name != "*"
}
