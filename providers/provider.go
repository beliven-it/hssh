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
	iGetEntity
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

type iGetEntity interface {
	GetEntity() string
}

/*
Provider use two attributes
url and privateToken.

url is the repo link where files can be found.

privateToken instead permit to
authenticate to the service
*/
type provider struct {
	url              string
	privateToken     string
	connectionString string
	driver           string
	entity           string
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

func (p *provider) GetEntity() string {
	return p.entity
}

func (p *provider) GetConnectionString() string {
	return p.connectionString
}

func (p *provider) GetPrivateToken() string {
	return p.privateToken
}

func (p *provider) GetURL() string {
	return p.url
}

func getDriverFromConnectionString(connectionString string) (string, error) {
	rgx := regexp.MustCompile("^(.*?)://")
	driver := rgx.FindAllStringSubmatch(connectionString, 1)

	if len(driver) == 0 {
		return "", errors.New("Invalid connection string")
	}

	return driver[0][1], nil
}

// New ...
/*............................................................................*/
func New(connectionString string) (IProvider, error) {
	driver, err := getDriverFromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	switch driver {
	case "gitlab":
		return NewGitlab(connectionString)
	case "github":
		return NewGithub(connectionString)
	default:
		return nil, errors.New("Invalid driver provider " + driver)
	}
}
