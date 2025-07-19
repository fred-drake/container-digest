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
	client *regclient.RegClient
}

// NewClient creates a new registry client
func NewClient(containersConfig *models.ContainersConfig) (*Client, error) {
	// Initialize regclient with Docker config
	rc := regclient.New(regclient.WithDockerCreds(), regclient.WithDockerCerts())
	
	client := &Client{
		client: rc,
	}
	
	return client, nil
}

// GetDigests fetches digests for all containers in the config
func (c *Client) GetDigests(containersConfig *models.ContainersConfig) (models.NestedDigestResults, error) {
	results := models.NestedDigestResults{}
	ctx := context.Background()

	for _, container := range containersConfig.Containers {
		// For each architecture, get the digest
		for _, arch := range container.Architectures {
			// Get the digest for this specific architecture
			digest, err := c.GetDigest(ctx, container.Repository, container.Name, container.Tag, arch)
			if err != nil {
				return nil, fmt.Errorf("failed to get digest for %s/%s:%s (%s): %w", 
					container.Repository, container.Name, container.Tag, arch, err)
			}

			// Initialize maps if they don't exist
			if _, exists := results[container.Repository]; !exists {
				results[container.Repository] = models.RepositoryMap{}
			}
			
			if _, exists := results[container.Repository][container.Name]; !exists {
				results[container.Repository][container.Name] = models.TagMap{}
			}
			
			if _, exists := results[container.Repository][container.Name][container.Tag]; !exists {
				results[container.Repository][container.Name][container.Tag] = models.ArchMap{}
			}

			// Add the digest to the nested structure
			results[container.Repository][container.Name][container.Tag][arch] = digest
		}
	}

	return results, nil
}

// GetDigest fetches the digest for a specific container and architecture
func (c *Client) GetDigest(ctx context.Context, registry, name, tag, architecture string) (string, error) {
	// Create the full reference string (registry/repository:tag)
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
func (c *Client) DebugManifest(registry, name, tag string) error {
	ctx := context.Background()
	
	// Create the full reference string (registry/repository:tag)
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
