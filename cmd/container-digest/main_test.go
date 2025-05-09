package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/fdrake/container-digest/internal/models"
)

// TestJSONKeysOrdering tests that JSON keys are ordered alphabetically
func TestJSONKeysOrdering(t *testing.T) {
	// Create test data with intentionally unordered keys
	testData := models.NestedDigestResults{
		"registry2": models.RepositoryMap{
			"repo2": models.TagMap{
				"tag2": models.ArchMap{
					"arch2": "digest2",
					"arch1": "digest1",
				},
				"tag1": models.ArchMap{
					"arch1": "digest1",
				},
			},
			"repo1": models.TagMap{},
		},
		"registry1": models.RepositoryMap{},
	}

	// Marshal using our custom JSON marshaler
	sortedJSON := &orderedJSON{Value: testData}
	jsonBytes, err := sortedJSON.MarshalIndent()
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	jsonString := string(jsonBytes)

	// Verify registry keys are ordered
	registry1Pos := strings.Index(jsonString, "registry1")
	registry2Pos := strings.Index(jsonString, "registry2")
	if registry1Pos > registry2Pos {
		t.Errorf("registry1 should appear before registry2 in alphabetical order")
	}

	// Verify repository keys are ordered
	repo1Pos := strings.Index(jsonString, "repo1")
	repo2Pos := strings.Index(jsonString, "repo2")
	if repo1Pos > repo2Pos {
		t.Errorf("repo1 should appear before repo2 in alphabetical order")
	}

	// Verify tag keys are ordered
	tag1Pos := strings.Index(jsonString, "tag1")
	tag2Pos := strings.Index(jsonString, "tag2")
	if tag1Pos > tag2Pos {
		t.Errorf("tag1 should appear before tag2 in alphabetical order")
	}

	// Verify arch keys are ordered
	arch1Pos := strings.Index(jsonString, "arch1")
	arch2Pos := strings.Index(jsonString, "arch2")
	if arch1Pos > arch2Pos {
		t.Errorf("arch1 should appear before arch2 in alphabetical order")
	}

	// Parse the JSON back to verify structure is preserved
	var parsedData models.NestedDigestResults
	err = json.Unmarshal(jsonBytes, &parsedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify the data structure is preserved
	if !reflect.DeepEqual(testData, parsedData) {
		t.Errorf("Original and parsed data structures don't match")
	}
}

// TestNixKeysOrdering tests that Nix keys are ordered alphabetically
func TestNixKeysOrdering(t *testing.T) {
	// Create test data with intentionally unordered keys
	testData := models.NestedDigestResults{
		"registry2": models.RepositoryMap{
			"repo2": models.TagMap{
				"tag2": models.ArchMap{
					"arch2": "digest2",
					"arch1": "digest1",
				},
				"tag1": models.ArchMap{
					"arch1": "digest1",
				},
			},
			"repo1": models.TagMap{},
		},
		"registry1": models.RepositoryMap{},
	}

	// Convert to Nix format
	nixOutput, err := formatAsNix(testData)
	if err != nil {
		t.Fatalf("Failed to format as Nix: %v", err)
	}

	// Verify registry keys are ordered
	registry1Pos := strings.Index(nixOutput, "\"registry1\"")
	registry2Pos := strings.Index(nixOutput, "\"registry2\"")
	if registry1Pos > registry2Pos {
		t.Errorf("registry1 should appear before registry2 in alphabetical order")
	}

	// Verify repository keys are ordered
	repo1Pos := strings.Index(nixOutput, "\"repo1\"")
	repo2Pos := strings.Index(nixOutput, "\"repo2\"")
	if repo1Pos > repo2Pos {
		t.Errorf("repo1 should appear before repo2 in alphabetical order")
	}

	// Verify tag keys are ordered
	tag1Pos := strings.Index(nixOutput, "\"tag1\"")
	tag2Pos := strings.Index(nixOutput, "\"tag2\"")
	if tag1Pos > tag2Pos {
		t.Errorf("tag1 should appear before tag2 in alphabetical order")
	}

	// Verify arch keys are ordered
	arch1Pos := strings.Index(nixOutput, "\"arch1\"")
	arch2Pos := strings.Index(nixOutput, "\"arch2\"")
	if arch1Pos > arch2Pos {
		t.Errorf("arch1 should appear before arch2 in alphabetical order")
	}
}

// TestGetSortedKeys tests the getSortedKeys helper function
func TestGetSortedKeys(t *testing.T) {
	// Test with a map[string]interface{}
	testMap := map[string]interface{}{
		"c": 3,
		"a": 1,
		"b": 2,
	}

	keys := getSortedKeys(testMap)
	
	// Check length
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Check order
	expectedKeys := []string{"a", "b", "c"}
	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Errorf("Expected key at position %d to be %s, got %s", i, expectedKeys[i], key)
		}
	}

	// Test with a non-map value
	nonMapKeys := getSortedKeys("not a map")
	if nonMapKeys != nil {
		t.Errorf("Expected nil for non-map input, got %v", nonMapKeys)
	}
}
