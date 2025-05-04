package models

// ContainersConfig represents the structure of containers.toml file
type ContainersConfig struct {
	Containers []Container `toml:"containers"` // List of containers to fetch digests for
}

// Container represents a container entry in the containers.toml file
type Container struct {
	Repository   string   `toml:"repository"`   // Repository hostname (e.g., docker.io, ghcr.io)
	Name         string   `toml:"name"`         // Container name (e.g., library/busybox)
	Tag          string   `toml:"tag"`          // Container tag (e.g., latest)
	Architectures []string `toml:"architectures"` // List of architectures (e.g., "linux/amd64", "linux/arm/v5")
}

// AuthConfig represents the structure of authentication.toml file
type AuthConfig struct {
	Credentials map[string]Credential `toml:"credentials"` // Key is the repository label, value is the credential
}

// Credential holds the username and password for a repository
type Credential struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// DigestResult represents a single container digest result
type DigestResult struct {
	Repository   string   `json:"repository"`
	Name         string   `json:"name"`
	Tag          string   `json:"tag"`
	Architectures []ArchDigest `json:"architectures"`
}

// ArchDigest represents the digest for a specific architecture
type ArchDigest struct {
	Architecture string `json:"architecture"`
	Digest       string `json:"digest"`
}

// DigestResults is a slice of DigestResult
type DigestResults []DigestResult
