package config

import (
	"fmt"

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
