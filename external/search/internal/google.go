package internal

// HTTP request fields and URI
//
// Defined https://developers.google.com/custom-search/json-api/v1/reference/cse/list#request
const (
	GoogleURI = "https://www.googleapis.com/customsearch/v1"

	RequestKey            = "key"
	RequestSearchEngineID = "cx"
	RequestQuery          = "q"
)

// SearchRequest is part of the Google Search API request.
//
// Note: this is public so that it can be used in the search package but "internal" so that it doesn't leak from the
// "external/search" package
type SearchRequest struct {
	// API key
	Key string

	// Custom Search Engine ID
	SearchEngineID string

	// query string
	Query string
}

// SearchResults is part of the Google Search API response
//
// Note: this is public so that the JSON parser works properly but "internal" so that it doesn't leak from the
// "external/search" package
//
// Defined https://developers.google.com/custom-search/json-api/v1/reference/cse/list#response
type SearchResults struct {
	Results []Result `json:"items"`
}

// Result is part of the Google Search API response
//
// Note: this is public so that the JSON parser works properly but "internal" so that it doesn't leak from the
// "external/search" package
//
// Defined https://developers.google.com/custom-search/json-api/v1/reference/cse/list#response
type Result struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}
