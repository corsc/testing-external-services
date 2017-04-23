package search

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/corsc/testing-external-services/external/search/internal"
	"github.com/stretchr/testify/assert"
)

// These tests verify our API contract with our users and therefore are UAT tests.
//
// They are designed to ensure that the code of this package is doing what we want it to independent of the external
// service.
//
// We do not consider them E2E tests as they do not make calls to the external service.

// Ensures that if there is no search term then we return ErrNoSearchTerm.
func TestApiImpl_Search_uat_emptyQuery(t *testing.T) {
	// define some test data
	searchEngineID := "aaa"
	key := "bbb"

	// define our expectations
	expected := Results{}
	expectedErr := ErrNoSearchTerm

	// build a client and send the request
	api := NewAPI(key, searchEngineID)
	result, resultErr := api.Search("")

	assert.Equal(t, expected, result)
	assert.Equal(t, expectedErr, resultErr)
}

// Happy path test - ensures we call the external API in the way we think we are calling it
func TestApiImpl_Search_uat_requestFormatValidation(t *testing.T) {
	// define some test data
	searchEngineID := "aaa"
	key := "bbb"
	query := "ccc"

	// build a test server to replace Google
	testServerWasCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		testServerWasCalled = true

		// check the request variables
		assert.Nil(t, req.ParseForm())
		assert.Equal(t, key, req.Form.Get(internal.RequestKey))
		assert.Equal(t, searchEngineID, req.Form.Get(internal.RequestSearchEngineID))
		assert.Equal(t, query, req.Form.Get(internal.RequestQuery))

		// return a valid but empty response
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(`{}`))
	}))

	// build a client and inject the test server URI
	api := NewAPI(key, searchEngineID)
	api.(*apiImpl).serviceURI = testServer.URL

	// send the request
	result, resultErr := api.Search(query)

	// validate the result
	assert.NotNil(t, result)
	assert.Nil(t, resultErr)
	assert.True(t, testServerWasCalled)
}

// Happy path test - ensures properly parse and convert the response
func TestApiImpl_Search_uat_responseParsing(t *testing.T) {
	// define some test data
	responseFromGoogle := `
{
	"items": [
		{
			"title": "aaa",
			"link": "http://aaa.com/",
			"snippet": "snip aaa"
		},
		{
			"title": "bbb",
			"link": "http://bbb.com/",
			"snippet": "snip bbb"
		},
		{
			"title": "ccc",
			"link": "http://ccc.com/",
			"snippet": "snip ccc"
		}
	]
}`

	expectedResults := Results{
		Results: []Result{
			{
				Title:   "aaa",
				Link:    "http://aaa.com/",
				Snippet: "snip aaa",
			},
			{
				Title:   "bbb",
				Link:    "http://bbb.com/",
				Snippet: "snip bbb",
			},
			{
				Title:   "ccc",
				Link:    "http://ccc.com/",
				Snippet: "snip ccc",
			},
		},
	}

	// build a test server to replace Google
	testServer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// return a valid response
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(responseFromGoogle))
	}))

	// build a client and inject the test server URI
	api := NewAPI("key", "engineID")
	api.(*apiImpl).serviceURI = testServer.URL

	// send the request
	results, resultErr := api.Search("apples")

	// validate the result
	assert.Equal(t, expectedResults, results)
	assert.Nil(t, resultErr)
}

// The following tests how the code in this package responds to errors from the external service.
// It is much easier to test against a mock service like this than to "spin up" the external service; particularly
// when it comes to testing error handling.
//
// These test will also provide contrast with the E2E tests.
// This is important because if these tests are working and the E2E tests break, then this most likely means the
// external service is no longer doing what you need it too (i.e. their API contract changed) or the configuration is
// incorrect.
//
// Personally I find this results in faster identification/debugging of issues.

func TestApiImpl_Search_uat_http400(t *testing.T) {
	// define some test data
	searchEngineID := "aaa"
	key := "bbb"
	query := "ccc"

	// build a test server to replace Google
	testServerWasCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		testServerWasCalled = true

		// return HTTP 400 error
		resp.WriteHeader(http.StatusBadRequest)
	}))

	// build a client and inject the test server URI
	api := NewAPI(key, searchEngineID)
	api.(*apiImpl).serviceURI = testServer.URL

	// send the request
	result, resultErr := api.Search(query)

	// validate the error is handled as we expect
	assert.Equal(t, EmptyResults, result)
	assert.NotNil(t, resultErr)
	assert.True(t, testServerWasCalled)
}
