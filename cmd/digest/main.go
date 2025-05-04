package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fdrake/container-digest/internal/config"
	"github.com/fdrake/container-digest/internal/registry"
)

func main() {
	// Define command-line flags
	containersFile := flag.String("containers", "containers.toml", "Path to containers TOML file")
	outputFile := flag.String("output", "", "Path to output JSON file (if not specified, output to stdout)")
	flag.Parse()

	// Load containers configuration
	containersConfig, err := config.LoadContainersConfig(*containersFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading containers config: %v\n", err)
		os.Exit(1)
	}

	// Create registry client
	client, err := registry.NewClient(containersConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating registry client: %v\n", err)
		os.Exit(1)
	}

	// Get digests for all containers
	results, err := client.GetDigests(containersConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching container digests: %v\n", err)
		os.Exit(1)
	}

	// Convert results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding results to JSON: %v\n", err)
		os.Exit(1)
	}

	// Output JSON data
	if *outputFile == "" {
		// Output to stdout
		fmt.Println(string(jsonData))
	} else {
		// Create parent directories if they don't exist
		if dir := filepath.Dir(*outputFile); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
				os.Exit(1)
			}
		}

		// Write to file
		if err := os.WriteFile(*outputFile, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Output written to %s\n", *outputFile)
	}
}
