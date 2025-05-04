package registry

import (
	"context"
	"fmt"
	"strings"

	"github.com/fdrake/container-digest/internal/models"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/platform"
	"github.com/regclient/regclient/types/ref"
)

// Client wraps the Docker registry client
type Client struct {
	client         *regclient.RegClient
	repoURLs       map[string]string
	repoCredentials map[string]models.Credential
}

// NewClient creates a new registry client
func NewClient(containersConfig *models.ContainersConfig, authConfig *models.AuthConfig) (*Client, error) {
	// Initialize regclient
	rc := regclient.New()
	
	// Store repository URLs and credentials for later use
	repoURLs := make(map[string]string)
	repoCredentials := make(map[string]models.Credential)
	
	// Copy repository URLs
	for repoKey, repoURL := range containersConfig.Repositories {
		repoURLs[repoKey] = repoURL
	}
	
	// Copy credentials if available
	if authConfig != nil {
		for repoKey, cred := range authConfig.Credentials {
			repoCredentials[repoKey] = cred
		}
	}
	
	client := &Client{
		client:          rc,
		repoURLs:        repoURLs,
		repoCredentials: repoCredentials,
	}
	
	return client, nil
}

// GetDigests fetches digests for all containers in the config
func (c *Client) GetDigests(containersConfig *models.ContainersConfig) (models.DigestResults, error) {
	results := models.DigestResults{}
	ctx := context.Background()

	for _, container := range containersConfig.Containers {
		// Check if the repository exists in the configuration
		_, ok := c.repoURLs[container.Repository]
		if !ok {
			return nil, fmt.Errorf("repository %s not found in configuration", container.Repository)
		}

		result := models.DigestResult{
			Repository:    container.Repository, // Use the repository key instead of the URL
			Name:          container.Name,
			Tag:           container.Tag,
			Architectures: []models.ArchDigest{},
		}

		// For each architecture, get the digest
		for _, arch := range container.Architectures {
			// Get the digest for this specific architecture
			digest, err := c.GetDigest(ctx, container.Repository, container.Name, container.Tag, arch)
			if err != nil {
				return nil, fmt.Errorf("failed to get digest for %s/%s:%s (%s): %w", 
					container.Repository, container.Name, container.Tag, arch, err)
			}

			// Add the digest to the result
			result.Architectures = append(result.Architectures, models.ArchDigest{
				Architecture: arch,
				Digest:       digest,
			})
		}

		results = append(results, result)
	}

	return results, nil
}

// GetDigest fetches the digest for a specific container and architecture
func (c *Client) GetDigest(ctx context.Context, repoKey, name, tag, architecture string) (string, error) {
	// Get the repository URL
	repoURL, ok := c.repoURLs[repoKey]
	if !ok {
		return "", fmt.Errorf("repository %s not found in configuration", repoKey)
	}

	// Create the full reference string (registry/repository:tag)
	// Extract registry domain from URL
	registry := strings.TrimPrefix(repoURL, "https://")
	registry = strings.TrimPrefix(registry, "http://")
	// Remove trailing slash if present
	registry = strings.TrimSuffix(registry, "/")
	
	// Create the full reference
	fullRef := fmt.Sprintf("%s/%s:%s", registry, name, tag)
	
	// Create image reference
	imageRef, err := ref.New(fullRef)
	if err != nil {
		return "", fmt.Errorf("failed to create image reference for %s: %w", fullRef, err)
	}

	// Parse the architecture string (e.g., "linux/amd64" -> OS: "linux", Architecture: "amd64")
	parts := strings.Split(architecture, "/")
	
	plat := platform.Platform{
		OS:           "linux", // Default to linux
		Architecture: "",
	}
	
	if len(parts) >= 2 {
		plat.OS = parts[0]
		plat.Architecture = parts[1]
		
		// Handle arm variants (e.g., "linux/arm/v7")
		if len(parts) >= 3 && parts[1] == "arm" {
			plat.Variant = parts[2]
		}
	} else {
		// If the format is not as expected, use the whole string as architecture
		plat.Architecture = architecture
	}
	
	// If architecture is still empty, default to amd64
	if plat.Architecture == "" {
		plat.Architecture = "amd64"
	}

	// First get the general manifest
	manifest, err := c.client.ManifestGet(ctx, imageRef)
	if err != nil {
		return "", fmt.Errorf("failed to get manifest for %s: %w", fullRef, err)
	}

	// If this is a manifest list (multi-arch), try to find the specific platform
	if manifest.IsList() {
		// Get the platform-specific descriptor
		platDesc, err := manifest.GetPlatformDesc(&plat)
		if err == nil && platDesc != nil {
			// We found a platform-specific manifest, return its digest
			return platDesc.Digest.String(), nil
		}
	}

	// Return the digest from the manifest (either single arch or couldn't find platform-specific)
	return manifest.GetDescriptor().Digest.String(), nil
}

// DebugManifest prints detailed information about a container manifest
func (c *Client) DebugManifest(repoKey, name, tag string) error {
	ctx := context.Background()
	
	// Get the repository URL
	repoURL, ok := c.repoURLs[repoKey]
	if !ok {
		return fmt.Errorf("repository %s not found in configuration", repoKey)
	}

	// Create the full reference string (registry/repository:tag)
	// Extract registry domain from URL
	registry := strings.TrimPrefix(repoURL, "https://")
	registry = strings.TrimPrefix(registry, "http://")
	// Remove trailing slash if present
	registry = strings.TrimSuffix(registry, "/")
	
	// Create the full reference
	fullRef := fmt.Sprintf("%s/%s:%s", registry, name, tag)
	
	// Create image reference
	imageRef, err := ref.New(fullRef)
	if err != nil {
		return fmt.Errorf("failed to create image reference for %s: %w", fullRef, err)
	}

	// Get manifest
	manifest, err := c.client.ManifestGet(ctx, imageRef)
	if err != nil {
		return fmt.Errorf("failed to get manifest for %s: %w", fullRef, err)
	}

	// Print manifest details
	fmt.Printf("Manifest Type: %s\n", manifest.GetMediaType())
	fmt.Printf("Manifest Digest: %s\n", manifest.GetDescriptor().Digest.String())
	
	// Check if this is a manifest list
	if manifest.IsList() {
		// For manifest lists, try to get platform list using the manifest API
		platformList, err := manifest.GetPlatformList()
		if err == nil && len(platformList) > 0 {
			fmt.Println("Available Platforms:")
			for _, plat := range platformList {
				fmt.Printf("  - %s/%s", plat.OS, plat.Architecture)
				if plat.Variant != "" {
					fmt.Printf("/%s", plat.Variant)
				}
				fmt.Println()
			}
		}
	}
	
	return nil
}
