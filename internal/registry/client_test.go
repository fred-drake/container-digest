package registry

import (
	"testing"

	"github.com/regclient/regclient"
)

// NewMockClient creates a mock client for testing
func NewMockClient() *Client {
	// We don't use real registry clients in tests
	return &Client{
		client: regclient.New(),
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
	
	// Verify the mock client was created successfully
}
