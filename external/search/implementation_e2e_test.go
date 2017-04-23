package search

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests are End to End (E2E) test on this package.
//
// They ensure that this package provides the functionality it intends to.
//
// They make "real" calls to the Google API as such they require and API key and custom search engine id to be
// supplied as environment variables ("KEY" and "CX" respectively)
//

func TestApiImpl_Search_e2e_happyPath(t *testing.T) {
	// pull the credentials from env
	key := os.Getenv(EnvAPIKey)
	if key == "" {
		t.Skip(fmt.Sprintf("Env var '%s' is required for this test", EnvAPIKey))
	}
	searchEngineID := os.Getenv(EnvSearchEngineID)
	if searchEngineID == "" {
		t.Skip(fmt.Sprintf("Env var '%s' is required for this test", EnvSearchEngineID))
	}

	client := NewAPI(key, searchEngineID)
	results, err := client.Search("mysql")
	assert.Nil(t, err)
	assert.True(t, len(results.Results) > 0)
}

func TestApiImpl_Search_e2e_badKey(t *testing.T) {
	// invalid key
	key := "invalid"

	// pull search engine ID from env
	searchEngineID := os.Getenv(EnvSearchEngineID)
	if searchEngineID == "" {
		t.Skip(fmt.Sprintf("Env var '%s' is required for this test", EnvSearchEngineID))
	}

	client := NewAPI(key, searchEngineID)
	results, err := client.Search("mysql")
	assert.Equal(t, ErrInvalidConfig, err)
	assert.True(t, len(results.Results) == 0)
}

func TestApiImpl_Search_e2e_badSearchEngineID(t *testing.T) {
	// pull the key from env
	key := os.Getenv(EnvAPIKey)
	if key == "" {
		t.Skip(fmt.Sprintf("Env var '%s' is required for this test", EnvAPIKey))
	}

	// invalid search engine ID
	searchEngineID := "invalid"

	client := NewAPI(key, searchEngineID)
	results, err := client.Search("mysql")
	assert.Equal(t, ErrInvalidConfig, err)
	assert.True(t, len(results.Results) == 0)
}
