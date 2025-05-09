package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
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

		// Convert results to JSON with alphabetically sorted keys
		encoder := json.NewEncoder(nil)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		sortedJSON := &orderedJSON{Value: transformedResults}
		outputData, err = sortedJSON.MarshalIndent()
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

// formatAsNix converts the digest results to Nix format with alphabetically sorted keys
func formatAsNix(results models.NestedDigestResults) (string, error) {
	var nixOutput string
	nixOutput = "{pkgs, ...}: {\n"

	// Get sorted registry keys
	registryKeys := getSortedKeys(results)
	
	// Iterate through registries in sorted order
	for _, registry := range registryKeys {
		repositories := results[registry]
		nixOutput += fmt.Sprintf("  \"%s\" = {\n", escapeNixString(registry))

		// Get sorted repository keys
		repoKeys := getSortedKeys(repositories)
		
		// Iterate through repositories in sorted order
		for _, repo := range repoKeys {
			tags := repositories[repo]
			nixOutput += fmt.Sprintf("    \"%s\" = {\n", escapeNixString(repo))

			// Get sorted tag keys
			tagKeys := getSortedKeys(tags)
			
			// Iterate through tags in sorted order
			for _, tag := range tagKeys {
				archs := tags[tag]
				nixOutput += fmt.Sprintf("      \"%s\" = {\n", escapeNixString(tag))

				// Get sorted architecture keys
				archKeys := getSortedKeys(archs)
				
				// Iterate through architectures in sorted order
				for _, arch := range archKeys {
					fullImageRef := archs[arch]
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

// getSortedKeys returns a sorted slice of map keys
func getSortedKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		return nil
	}
	
	keys := make([]string, 0, v.Len())
	for _, k := range v.MapKeys() {
		keys = append(keys, k.String())
	}
	
	sort.Strings(keys)
	return keys
}

// orderedJSON is a wrapper for marshaling JSON with ordered keys
type orderedJSON struct {
	Value interface{}
}

// MarshalIndent returns JSON with alphabetically sorted keys
func (o *orderedJSON) MarshalIndent() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	
	if err := o.marshalValue(buffer, o.Value, 0); err != nil {
		return nil, err
	}
	
	return buffer.Bytes(), nil
}

// marshalValue recursively marshals a value with ordered keys
func (o *orderedJSON) marshalValue(buffer *bytes.Buffer, v interface{}, indent int) error {
	if v == nil {
		buffer.WriteString("null")
		return nil
	}
	
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		return o.marshalMap(buffer, rv, indent)
	case reflect.Slice, reflect.Array:
		return o.marshalSlice(buffer, rv, indent)
	default:
		// For non-container types, use standard JSON marshaling
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		buffer.Write(b)
		return nil
	}
}

// marshalMap marshals a map with alphabetically sorted keys
func (o *orderedJSON) marshalMap(buffer *bytes.Buffer, rv reflect.Value, indent int) error {
	buffer.WriteString("{")
	
	// Get and sort the keys
	keys := make([]string, 0, rv.Len())
	for _, k := range rv.MapKeys() {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)
	
	// Write each key-value pair
	for i, key := range keys {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n" + strings.Repeat("  ", indent+1))
		
		// Marshal the key
		keyJSON, err := json.Marshal(key)
		if err != nil {
			return err
		}
		buffer.Write(keyJSON)
		
		buffer.WriteString(": ")
		
		// Marshal the value
		value := rv.MapIndex(reflect.ValueOf(key))
		if err := o.marshalValue(buffer, value.Interface(), indent+1); err != nil {
			return err
		}
	}
	
	if len(keys) > 0 {
		buffer.WriteString("\n" + strings.Repeat("  ", indent))
	}
	buffer.WriteString("}")
	return nil
}

// marshalSlice marshals a slice or array
func (o *orderedJSON) marshalSlice(buffer *bytes.Buffer, rv reflect.Value, indent int) error {
	buffer.WriteString("[")
	
	// Write each element
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n" + strings.Repeat("  ", indent+1))
		
		// Marshal the element
		if err := o.marshalValue(buffer, rv.Index(i).Interface(), indent+1); err != nil {
			return err
		}
	}
	
	if rv.Len() > 0 {
		buffer.WriteString("\n" + strings.Repeat("  ", indent))
	}
	buffer.WriteString("]")
	return nil
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
