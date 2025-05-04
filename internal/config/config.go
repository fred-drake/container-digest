package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/fdrake/container-digest/internal/models"
)

// LoadContainersConfig loads container configuration from a TOML file
func LoadContainersConfig(path string) (*models.ContainersConfig, error) {
	config := &models.ContainersConfig{}
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, fmt.Errorf("failed to decode containers config: %w", err)
	}
	return config, nil
}

// LoadAuthConfig loads authentication configuration from a TOML file
// Returns nil config and nil error if the file doesn't exist
func LoadAuthConfig(path string) (*models.AuthConfig, error) {
	config := &models.AuthConfig{
		Credentials: make(map[string]models.Credential),
	}
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist, return empty config
		return config, nil
	}
	
	if _, err := toml.DecodeFile(path, config); err != nil {
		// If the file exists but can't be read, return an error
		return nil, fmt.Errorf("failed to decode auth config: %w", err)
	}
	
	return config, nil
}
