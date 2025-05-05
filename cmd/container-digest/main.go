package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fdrake/container-digest/internal/config"
	"github.com/fdrake/container-digest/internal/models"
	"github.com/fdrake/container-digest/internal/registry"
	"github.com/spf13/cobra"
)

var (
	containersFile string
	outputFile     string
	outputFormat   string
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

	// Generate output based on format
	var outputData []byte
	var formatName string

	switch outputFormat {
	case "json":
		// Transform results to include full image references
		transformedResults, err := transformResultsWithFullRefs(results)
		if err != nil {
			return fmt.Errorf("error transforming results: %w", err)
		}
		
		// Convert results to JSON
		outputData, err = json.MarshalIndent(transformedResults, "", "  ")
		if err != nil {
			return fmt.Errorf("error encoding results to JSON: %w", err)
		}
		formatName = "JSON"
	case "nix":
		// Transform results to include full image references
		transformedResults, err := transformResultsWithFullRefs(results)
		if err != nil {
			return fmt.Errorf("error transforming results: %w", err)
		}
		
		// Convert results to Nix format
		nixOutput, err := formatAsNix(transformedResults)
		if err != nil {
			return fmt.Errorf("error encoding results to Nix format: %w", err)
		}
		outputData = []byte(nixOutput)
		formatName = "Nix"
	default:
		return fmt.Errorf("unsupported output format: %s (supported formats: json, nix)", outputFormat)
	}

	// Output data
	if outputFile == "" {
		// Output to stdout
		fmt.Println(string(outputData))
	} else {
		// Create parent directories if they don't exist
		if dir := filepath.Dir(outputFile); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("error creating output directory: %w", err)
			}
		}

		// Write to file
		if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
			return fmt.Errorf("error writing output to file: %w", err)
		}
		fmt.Printf("%s output written to %s\n", formatName, outputFile)
	}

	return nil
}

// formatAsNix converts the digest results to Nix format
func formatAsNix(results models.NestedDigestResults) (string, error) {
	var nixOutput string
	nixOutput = "{\n"

	// Iterate through registries
	for registry, repositories := range results {
		nixOutput += fmt.Sprintf("  \"%s\" = {\n", escapeNixString(registry))

		// Iterate through repositories
		for repo, tags := range repositories {
			nixOutput += fmt.Sprintf("    \"%s\" = {\n", escapeNixString(repo))

			// Iterate through tags
			for tag, archs := range tags {
				nixOutput += fmt.Sprintf("      \"%s\" = {\n", escapeNixString(tag))

				// Iterate through architectures
				for arch, fullImageRef := range archs {
					nixOutput += fmt.Sprintf("        \"%s\" = \"%s\";\n",
						escapeNixString(arch), escapeNixString(fullImageRef))
				}

				nixOutput += "      };\n"
			}

			nixOutput += "    };\n"
		}

		nixOutput += "  };\n"
	}

	nixOutput += "}"
	return nixOutput, nil
}

// escapeNixString escapes special characters in strings for Nix format
func escapeNixString(s string) string {
	// Replace any special characters as needed
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "$", "\\$")
	return s
}

// transformResultsWithFullRefs transforms the nested digest results to include full image references
func transformResultsWithFullRefs(results models.NestedDigestResults) (models.NestedDigestResults, error) {
	transformedResults := models.NestedDigestResults{}
	
	// Iterate through registries
	for registry, repositories := range results {
		transformedResults[registry] = models.RepositoryMap{}
		
		// Iterate through repositories
		for repo, tags := range repositories {
			transformedResults[registry][repo] = models.TagMap{}
			
			// Iterate through tags
			for tag, archs := range tags {
				transformedResults[registry][repo][tag] = models.ArchMap{}
				
				// Iterate through architectures
				for arch, digest := range archs {
					// Format the full image reference with digest
					fullImageRef := fmt.Sprintf("%s/%s@%s", registry, repo, digest)
					transformedResults[registry][repo][tag][arch] = fullImageRef
				}
			}
		}
	}
	
	return transformedResults, nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "container-digest",
		Short: "Get container image digests from registries",
		Long:  `container-digest reads a TOML file containing docker container information and returns the sha256 digests of those containers, along with tags and architectures.`,
		RunE:  runDigest,
	}

	// Define command-line flags
	rootCmd.Flags().StringVar(&containersFile, "containers", "containers.toml", "Path to containers TOML file")
	rootCmd.Flags().StringVar(&outputFile, "output", "", "Path to output file (if not specified, output to stdout)")
	rootCmd.Flags().StringVar(&outputFormat, "output-format", "json", "Output format (json or nix)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
