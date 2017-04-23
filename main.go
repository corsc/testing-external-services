// Package main is the main entry point for this "service".
package main

import (
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/corsc/testing-external-services/external/search"
)

func main() {
	// standard "no frills" HTTP server copied from https://golang.org/pkg/net/http/
	http.HandleFunc("/", myHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// This is total contrived and badly implemented
func myHandler(resp http.ResponseWriter, req *http.Request) {
	// extract the search ter,
	err := req.ParseForm()
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Bad request. " + err.Error()))
		return
	}

	term := req.Form.Get("q")
	if term == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Bad request. Request should be ?q=[term]"))
		return
	}

	results, err := getSearchClient().Search(term)
	if err != nil {
		switch err {
		case search.ErrInvalidConfig:
			resp.WriteHeader(http.StatusInternalServerError)
			resp.Write([]byte("Bad config. " + err.Error()))
			return

		default:
			resp.WriteHeader(http.StatusServiceUnavailable)
			resp.Write([]byte("Bad request. " + err.Error()))
			return
		}
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprintf("OK - %d records found\n", len(results.Results))))

	for _, result := range results.Results {
		resp.Write([]byte(result.Title + " - " + result.Link + "\n"))
	}
}

func getSearchClient() search.Client {
	return search.NewAPI(os.Getenv(search.EnvAPIKey), os.Getenv(search.EnvSearchEngineID))
}
