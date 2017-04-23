package search

import "errors"

var (
	// ErrInvalidConfig indicates the api was incorrectly configured (check logs for more information)
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrNoSearchTerm indicates no search term was supplied
	ErrNoSearchTerm = errors.New("no search term supplied")

	// empty result
	EmptyResults = Results{}
)

// Environmental variables used to configure this service
const (
	// Google Search API Key
	EnvAPIKey = "KEY"

	// Google Custom Search Engine ID
	EnvSearchEngineID = "CX"
)

// Client defines the exported API for this package
type Client interface {
	Search(term string) (Results, error)
}

// Results contains the results of a single query
//
// NOTE: I am intentionally not using the Google format to ensure loose coupling between the API exported from this
// package and the external service.  I am also intentionally hiding all HTTP/REST/JSON evidence.
// Should the external service, its API format or transport mechanisms change in the future we may be able to maintain
// a consistent API between this package it's users.
type Results struct {
	Results []Result
}

// Result contains an individual search result
type Result struct {
	Title   string
	Link    string
	Snippet string
}
