package registry

import (
	"testing"

	"github.com/fdrake/container-digest/internal/models"
	"github.com/regclient/regclient"
)

// NewMockClient creates a mock client for testing
func NewMockClient() *Client {
	// We don't use real registry clients in tests
	return &Client{
		client:          regclient.New(),
		repoURLs:        make(map[string]string),
		repoCredentials: make(map[string]models.Credential),
	}
}

func TestMockClient(t *testing.T) {
	// This is a simplified test that just checks the client constructor logic
	// without making actual registry calls
	mockClient := NewMockClient()
	if mockClient == nil {
		t.Fatal("Failed to create mock client")
	}
	
	// Verify the mock client was created
	if mockClient.client == nil {
		t.Error("Mock client has nil regclient")
	}
	
	// Verify the mock client has empty maps
	if len(mockClient.repoURLs) != 0 {
		t.Error("Mock client should have empty repoURLs map")
	}
	
	if len(mockClient.repoCredentials) != 0 {
		t.Error("Mock client should have empty repoCredentials map")
	}
}
