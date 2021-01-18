package providers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type gitlab struct {
	provider
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

func (g *gitlab) Start() *gitlab {
	_, err := g.provider.ParseConnection("gitlab")
	if err != nil {
		log.Fatal(err)
	}

	return g
}

// NewGitlab ...
/*............................................................................*/
func NewGitlab(connectionString string) IProvider {
	g := gitlab{
		provider: provider{
			connectionString: connectionString,
			url:              "https://gitlab.com/api/v4",
		},
	}
	return g.Start()
}
