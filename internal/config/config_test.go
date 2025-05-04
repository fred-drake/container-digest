package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadContainersConfig(t *testing.T) {
	// Create a temporary TOML file for testing
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test-containers.toml")

	// Write test TOML content
	tomlContent := `
# Test container configuration

[[containers]]
repository = "docker.io"
name = "library/busybox"
tag = "latest"
architectures = ["linux/amd64", "linux/arm/v5"]

[[containers]]
repository = "ghcr.io"
name = "user/repo"
tag = "1.0.0"
architectures = ["linux/amd64"]
`
	err := os.WriteFile(tmpFile, []byte(tomlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test loading the config
	config, err := LoadContainersConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadContainersConfig returned an error: %v", err)
	}

	// Verify the loaded config

	if len(config.Containers) != 2 {
		t.Errorf("Expected 2 containers, got %d", len(config.Containers))
	}

	if config.Containers[0].Repository != "docker.io" {
		t.Errorf("Expected first container repository to be 'docker.io', got '%s'", config.Containers[0].Repository)
	}

	if config.Containers[0].Name != "library/busybox" {
		t.Errorf("Expected first container name to be 'library/busybox', got '%s'", config.Containers[0].Name)
	}

	if len(config.Containers[0].Architectures) != 2 {
		t.Errorf("Expected 2 architectures for first container, got %d", len(config.Containers[0].Architectures))
	}
}
