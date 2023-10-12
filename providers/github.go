package providers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type githubFile struct {
	ID   string `json:"sha"`
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type github struct {
	provider
}

// get
/*............................................................................*/
func (g *github) get(endpoint string, queryParams []queryParam) ([]byte, error) {
	request, err := http.NewRequest("GET", g.url+endpoint, nil)
	if err != nil {
		return nil, err
	}

	if g.connection.Token != "" {
		request.Header.Set("Authorization", "token "+g.connection.Token)
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

	body, err := io.ReadAll(response.Body)
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
func (g *github) GetFiles(repo string, filePath string) ([]file, error) {
	endpoint := "/repos/" + repo + "/contents/" + filePath

	queryParams := []queryParam{}
	bodyInBytes, err := g.get(endpoint, queryParams)
	if err != nil {
		return nil, err
	}

	gitFiles := []githubFile{}
	f := []file{}

	err = json.Unmarshal(bodyInBytes, &gitFiles)
	if err != nil {
		return nil, err
	}

	for _, gf := range gitFiles {
		if gf.Type != "file" {
			continue
		}

		newFile := file{
			ID:   gf.ID,
			Name: gf.Name,
			Path: gf.Path,
		}

		f = append(f, newFile)
	}

	return f, nil
}

// GetFile
/*............................................................................*/
func (g *github) GetFile(repo string, fileID string) ([]byte, error) {
	var content []byte
	endpoint := "/repos/" + repo + "/git/blobs/" + fileID
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

// NewGithub ...
/*............................................................................*/
func NewGithub(connection ProviderConnection) (IProvider, error) {
	if connection.URL == "" {
		connection.URL = "https://api.github.com"
	}

	g := github{
		provider: provider{
			driver:     "github",
			connection: connection,
			url:        connection.URL,
		},
	}
	return &g, nil
}
