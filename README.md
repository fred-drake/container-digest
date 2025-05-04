# Container Digest

A Go application that reads a TOML file containing Docker container information and returns the SHA256 digests of those containers, along with their tags and architectures.

## Features

- Reads container information from a TOML configuration file
- Supports multiple architectures per container
- Outputs digests in JSON format

## Installation

```bash
go install github.com/fdrake/container-digest/cmd/digest
```

Or build from source:

```bash
git clone https://github.com/fdrake/container-digest.git
cd container-digest
go build -o digest ./cmd/digest
```

## Usage

```bash
./digest -containers=containers.toml -output=digests.json
```

### Command-line options

- `-containers`: Path to the containers TOML file (default: "containers.toml")
- `-output`: Path to the output JSON file (if not specified, output to stdout)

## Configuration

### Containers Configuration (containers.toml)

```toml
[repositories]
docker = "https://registry-1.docker.io"
github = "https://ghcr.io"

[[containers]]
repository = "docker"
name = "library/busybox"
tag = "latest"
architectures = ["linux/amd64", "linux/arm64"]

[[containers]]
repository = "github"
name = "user/repo"
tag = "main"
architectures = ["linux/amd64"]
```

## Output Format

The application outputs a JSON array of container digest information:

```json
[
  {
    "repository": "https://registry-1.docker.io",
    "name": "library/busybox",
    "tag": "latest",
    "architectures": [
      {
        "architecture": "linux/amd64",
        "digest": "sha256:abcdef..."
      },
      {
        "architecture": "linux/arm64",
        "digest": "sha256:123456..."
      }
    ]
  }
]
```

## Testing

Run the tests:

```bash
go test ./...
```

## Dependencies

- [github.com/heroku/docker-registry-client](https://github.com/heroku/docker-registry-client) - Docker Registry API client
- [github.com/BurntSushi/toml](https://github.com/BurntSushi/toml) - TOML parser
