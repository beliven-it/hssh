package providers

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type heimdallr struct {
	provider
}

// get
/*............................................................................*/
func (g *heimdallr) get(endpoint string, queryParams []queryParam) ([]byte, error) {
	request, err := http.NewRequest("GET", g.url+endpoint, nil)
	if err != nil {
		return nil, err
	}

	if g.privateToken != "" {
		request.Header.Set("Authorization", "Bearer "+g.privateToken)
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
func (g *heimdallr) GetFiles(projectID string, filePath string) ([]file, error) {
	f := []file{}

	mainFile := file{
		Name:    "hosts",
		Path:    "hosts",
		ID:      "hostnames",
		Content: "",
	}

	f = append(f, mainFile)

	return f, nil
}

// GetFile
/*............................................................................*/
func (g *heimdallr) GetFile(projectID string, fileID string) ([]byte, error) {
	var content []byte
	endpoint := "/hostnames"
	bodyInBytes, err := g.get(endpoint, []queryParam{})
	if err != nil {
		return content, err
	}

	return bodyInBytes, nil
}

func (g *heimdallr) Start() (*heimdallr, error) {
	_, err := g.provider.ParseConnection("heimdallr")
	if err != nil {
		return nil, err
	}

	return g, nil
}

// NewGitlab ...
/*............................................................................*/
func NewHeimdallr(connectionString string) (IProvider, error) {
	g := heimdallr{
		provider: provider{
			driver:           "heimdallr",
			connectionString: connectionString,
			url:              getOptionURLFromConnectionString(connectionString),
		},
	}

	return g.Start()
}
