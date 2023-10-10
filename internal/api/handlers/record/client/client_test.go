// BEGIN: 8f5d4a6b7c8d
package client

import (
	"testing"
)

func TestNewApiClient(t *testing.T) {
	endpoint := "https://example.com"
	bearerToken := "my-token"
	insecure := false

	// Test with valid inputs
	c, err := NewApiClient(endpoint, bearerToken, insecure)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if c == nil {
		t.Errorf("Expected a non-nil api.Client instance")
	}

	// Test with invalid bearer token
	_, err = NewApiClient(endpoint, "", insecure)
	if err == nil {
		t.Errorf("Expected an error when bearer token is empty")
	}

	// Test with invalid endpoint
	_, err = NewApiClient("", bearerToken, insecure)
	if err == nil {
		t.Errorf("Expected an error when endpoint is empty")
	}
}

// END: 8f5d4a6b7c8d
