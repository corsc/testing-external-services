package search

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"io/ioutil"

	"github.com/corsc/testing-external-services/external/search/internal"
	"github.com/myteksi/go/Godeps/_workspace/src/github.com/dustin/gojson"
)

// NewAPI returns a ready to use search API with the supplied credentials
//
// Note: Many people recommend not using constructors, exporting the implementation type and leaving the mocking to the
// user.  I do not.  I find that exporting only the API and providing mocks, I both force myself to use the interface
// (which results in better encapsulation) and providing the mocks reduces duplication.
func NewAPI(key string, searchEngineID string) Client {
	return &apiImpl{
		key:            key,
		searchEngineID: searchEngineID,
	}
}

// apiImpl implements the API
type apiImpl struct {
	key            string
	searchEngineID string
	serviceURI     string
}

// Search implements API
func (api *apiImpl) Search(term string) (Results, error) {
	req, err := api.buildRequest(term)
	if err != nil {
		return EmptyResults, err
	}

	resp, err := api.makeRequest(req)
	if err != nil {
		return EmptyResults, err
	}

	// decide what to do based on the response code
	switch resp.StatusCode {
	case http.StatusOK:
		return api.buildOutput(resp)

	case http.StatusBadRequest:
		return EmptyResults, ErrInvalidConfig

	default:
		return EmptyResults, fmt.Errorf("response code was unexpected. was %d ", resp.StatusCode)
	}

}

func (api *apiImpl) buildRequest(term string) (*internal.SearchRequest, error) {
	term = strings.TrimSpace(term)
	if term == "" {
		return nil, ErrNoSearchTerm
	}

	return &internal.SearchRequest{
		Key:            api.key,
		SearchEngineID: api.searchEngineID,
		Query:          term,
	}, nil
}

func (api *apiImpl) makeRequest(req *internal.SearchRequest) (*http.Response, error) {
	// build the request
	values := url.Values{}
	values.Add(internal.RequestKey, req.Key)
	values.Add(internal.RequestSearchEngineID, req.SearchEngineID)
	values.Add(internal.RequestQuery, req.Query)

	url := api.getServiceURI() + "?" + values.Encode()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//
func (api *apiImpl) getServiceURI() string {
	if api.serviceURI == "" {
		// default to the live endpoint
		return internal.GoogleURI
	}
	return api.serviceURI
}

func (api *apiImpl) buildOutput(resp *http.Response) (Results, error) {
	// parse Google's HTTP response
	serviceResponseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return EmptyResults, err
	}

	// parse into JSON
	serviceResponse := internal.SearchResults{}
	err = json.Unmarshal(serviceResponseBytes, &serviceResponse)
	if err != nil {
		return EmptyResults, err
	}

	// convert into our format
	results := Results{}
	for _, item := range serviceResponse.Results {
		result := Result{
			Title:   item.Title,
			Link:    item.Link,
			Snippet: item.Snippet,
		}

		results.Results = append(results.Results, result)
	}

	return results, nil
}
