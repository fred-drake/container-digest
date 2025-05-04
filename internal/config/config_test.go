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
[repositories]
docker = "https://docker.io"
github = "https://ghcr.io"

[[containers]]
repository = "docker"
name = "library/busybox"
tag = "latest"
architectures = ["linux/amd64", "linux/arm/v5"]

[[containers]]
repository = "github"
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
	if len(config.Repositories) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(config.Repositories))
	}

	if config.Repositories["docker"] != "https://docker.io" {
		t.Errorf("Expected docker repository URL to be 'https://docker.io', got '%s'", config.Repositories["docker"])
	}

	if len(config.Containers) != 2 {
		t.Errorf("Expected 2 containers, got %d", len(config.Containers))
	}

	if config.Containers[0].Name != "library/busybox" {
		t.Errorf("Expected first container name to be 'library/busybox', got '%s'", config.Containers[0].Name)
	}

	if len(config.Containers[0].Architectures) != 2 {
		t.Errorf("Expected 2 architectures for first container, got %d", len(config.Containers[0].Architectures))
	}
}

func TestLoadAuthConfig(t *testing.T) {
	// Create a temporary TOML file for testing
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test-auth.toml")

	// Write test TOML content
	tomlContent := `
[credentials.docker]
username = "user1"
password = "pass1"

[credentials.github]
username = "user2"
password = "pass2"
`
	err := os.WriteFile(tmpFile, []byte(tomlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test loading the config
	config, err := LoadAuthConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadAuthConfig returned an error: %v", err)
	}

	// Verify the loaded config
	if len(config.Credentials) != 2 {
		t.Errorf("Expected 2 credentials, got %d", len(config.Credentials))
	}

	if config.Credentials["docker"].Username != "user1" {
		t.Errorf("Expected docker username to be 'user1', got '%s'", config.Credentials["docker"].Username)
	}

	// Test non-existent file (should return empty config, not error)
	nonExistentFile := filepath.Join(tmpDir, "non-existent.toml")
	config, err = LoadAuthConfig(nonExistentFile)
	if err != nil {
		t.Errorf("LoadAuthConfig with non-existent file returned an error: %v", err)
	}
	if config == nil {
		t.Errorf("LoadAuthConfig with non-existent file returned nil config")
	}
	if len(config.Credentials) != 0 {
		t.Errorf("Expected empty credentials for non-existent file, got %d entries", len(config.Credentials))
	}
}
