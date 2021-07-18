package providers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

type gitlab struct {
	provider
	filePath string
	branch   string
}

// get
/*............................................................................*/
func (g *gitlab) get(endpoint string, queryParams []queryParam) ([]byte, error) {
	request, err := http.NewRequest("GET", g.url+endpoint, nil)
	if err != nil {
		return nil, err
	}

	if g.privateToken != "" {
		request.Header.Set("PRIVATE-TOKEN", g.privateToken)
	}

	query := request.URL.Query()

	for _, param := range queryParams {
		query.Add(param.key, param.value)
	}

	request.URL.RawQuery = query.Encode()

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, errors.New(string(body))
	}

	return body, nil
}

// TODO:
/*
Request return a limited set of results
if you need more than 20 files
you must create a paginatio system
*/
// GetFiles
/*............................................................................*/
func (g *gitlab) GetFiles(projectID string, filePath string) ([]file, error) {
	endpoint := "/projects/" + projectID + "/repository/tree"
	path := queryParam{
		key:   "path",
		value: filePath,
	}

	queryParams := []queryParam{path}
	bodyInBytes, err := g.get(endpoint, queryParams)
	if err != nil {
		return nil, err
	}

	f := []file{}
	err = json.Unmarshal(bodyInBytes, &f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// GetFile
/*............................................................................*/
func (g *gitlab) GetFile(projectID string, fileID string) ([]byte, error) {
	var content []byte
	endpoint := "/projects/" + projectID + "/repository/blobs/" + fileID
	bodyInBytes, err := g.get(endpoint, []queryParam{})
	if err != nil {
		return content, err
	}

	f := file{}
	err = json.Unmarshal(bodyInBytes, &f)
	if err != nil {
		return content, err
	}

	content, err = base64.StdEncoding.DecodeString(f.Content)
	if err != nil {
		return content, err
	}

	return content, nil
}

func (g *gitlab) parse() error {
	rgx := regexp.MustCompile(`(.*):\/\/(.*?)(:/|\s)(.*?)(@|\s)(.*?)(##|\s)(.*)`)
	result := rgx.FindAllStringSubmatch(g.connectionString, 1)

	if len(result) == 0 || len(result[0]) < 2 {
		return errors.New("Cannot extract token from connection string.\nThe connection string must follow the format:\n<driver>://<token>")
	}

	g.privateToken = result[0][2]
	g.entity = result[0][4]
	g.filePath = result[0][6]
	g.branch = result[0][8]

	return nil
}

// NewGitlab ...
/*............................................................................*/
func NewGitlab(connectionString string) (IProvider, error) {
	g := gitlab{
		provider: provider{
			driver:           "gitlab",
			connectionString: connectionString,
			url:              "https://gitlab.com/api/v4",
		},
	}
	return &g, g.parse()
}
