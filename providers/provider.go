package providers

import (
	"errors"
	"regexp"
)

/*
Provider is a abstract class that decscribe
the concrete classes used to fetch the connections
files from a remote repository
*/

// IProvider ...
type IProvider interface {
	iGet
	iGetFile
	iGetFiles
	iGetPrivateToken
	iGetDriver
}

type ProviderConnection struct {
	Type     string `mapstructure:"type"`
	URL      string `mapstructure:"url"`
	Token    string `mapstructure:"access_token"`
	EntityID string `mapstructure:"entity_id"`
	Subpath  string `mapstructure:"subpath"`
}

func (p *ProviderConnection) FromString(connectionString string) {
	rgx := regexp.MustCompile(`^(.*?)://(.*?):/(.*?)\@(.*)`)
	result := rgx.FindAllStringSubmatch(connectionString, 1)

	if len(result) >= 1 {
		p.Type = result[0][1]
		p.URL = ""
		p.Token = result[0][2]
		p.EntityID = result[0][3]
		p.Subpath = result[0][4]
	}
}

type iGet interface {
	get(string, []queryParam) ([]byte, error)
}

type iGetDriver interface {
	GetDriver() string
}

type iGetFiles interface {
	GetFiles(string, string) ([]file, error)
}

type iGetFile interface {
	GetFile(string, string) ([]byte, error)
}

type iGetPrivateToken interface {
	GetPrivateToken() string
}

/*
Provider use two attributes
url and privateToken.

url is the repo link where files can be found.

privateToken instead permit to
authenticate to the service
*/
type provider struct {
	url          string
	privateToken string
	connection   ProviderConnection
	driver       string
}

type file struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Name    string `json:"file_name"`
	Path    string `json:"path"`
}

type queryParam struct {
	key   string
	value string
}

func (p *provider) GetDriver() string {
	return p.driver
}

func (p *provider) GetPrivateToken() string {
	return p.privateToken
}

func (p *provider) GetURL() string {
	return p.url
}

// New ...
/*............................................................................*/
func New(connection ProviderConnection) (IProvider, error) {
	switch connection.Type {
	case "gitlab":
		return NewGitlab(connection)
	case "github":
		return NewGithub(connection)
	default:
		return nil, errors.New("Invalid driver provider " + connection.Type)
	}
}
