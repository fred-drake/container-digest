package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fdrake/container-digest/internal/config"
	"github.com/fdrake/container-digest/internal/registry"
	"github.com/spf13/cobra"
)

var (
	containersFile string
	outputFile     string
)

func runDigest(cmd *cobra.Command, args []string) error {
	// Load containers configuration
	containersConfig, err := config.LoadContainersConfig(containersFile)
	if err != nil {
		return fmt.Errorf("error loading containers config: %w", err)
	}

	// Create registry client
	client, err := registry.NewClient(containersConfig)
	if err != nil {
		return fmt.Errorf("error creating registry client: %w", err)
	}

	// Get digests for all containers
	results, err := client.GetDigests(containersConfig)
	if err != nil {
		return fmt.Errorf("error fetching container digests: %w", err)
	}

	// Convert results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding results to JSON: %w", err)
	}

	// Output JSON data
	if outputFile == "" {
		// Output to stdout
		fmt.Println(string(jsonData))
	} else {
		// Create parent directories if they don't exist
		if dir := filepath.Dir(outputFile); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("error creating output directory: %w", err)
			}
		}

		// Write to file
		if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
			return fmt.Errorf("error writing output to file: %w", err)
		}
		fmt.Printf("Output written to %s\n", outputFile)
	}

	return nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "container-digest",
		Short: "Get container image digests from registries",
		Long:  `container-digest reads a TOML file containing docker container information and returns a JSON of the sha256 digests of those containers, along with tags and architectures.`,
		RunE:  runDigest,
	}

	// Define command-line flags
	rootCmd.Flags().StringVar(&containersFile, "containers", "containers.toml", "Path to containers TOML file")
	rootCmd.Flags().StringVar(&outputFile, "output", "", "Path to output JSON file (if not specified, output to stdout)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
