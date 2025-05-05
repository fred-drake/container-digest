# Container Digest

A Go application that reads a TOML file containing Docker container information and returns the SHA256 digests of those containers, along with their tags and architectures.

I use this to declaratively configure my container images in my Nix configurations.

## Features

- Reads container information from a TOML configuration file
- Supports multiple architectures per container
- Outputs digests in JSON or Nix format

## Installation

The `devenv` Nix tool is used with `direnv` to pull the proper Go environment and supporting applications.

The `just` tool is used to build and run the application.

Available recipes:

- `build` # Build the container-digest application
- `clean` # Clean up the build directory
- `default` # List all available recipes with descriptions
- `run` # Run the application
- `test` # Run all unit tests

### Command-line options

- `--containers`: Path to the containers TOML file (default: "containers.toml")
- `--output`: Path to the output file (if not specified, output to stdout)
- `--output-format`: Output format, either "json" or "nix" (default: "json")

## Configuration

### Containers Configuration (containers.toml)

```toml
[[containers]]
repository = "docker.io"
name = "library/busybox"
tag = "latest"
architectures = ["linux/amd64", "linux/arm64", "linux/arm/v7"]

[[containers]]
repository = "docker.io"
name = "library/postgres"
tag = "16-alpine"
architectures = ["linux/amd64", "linux/arm64"]

[[containers]]
repository = "ghcr.io"
name = "home-assistant/home-assistant"
tag = "latest"
architectures = ["linux/amd64"]

[[containers]]
repository = "docker.gitea.com"
name = "gitea"
tag = "latest"
architectures = ["linux/amd64"]
```

## Output Formats

The application supports two output formats: JSON and Nix.

### JSON Format

When using the default `--output-format=json`, the application outputs a JSON structure of container digest information:

```json
{
  "docker.gitea.com": {
    "gitea": {
      "latest": {
        "linux/amd64": "sha256:5ee30f...de6367"
      }
    }
  },
  "docker.io": {
    "library/busybox": {
      "latest": {
        "linux/amd64": "sha256:ad9fa4...948f9f",
        "linux/arm/v7": "sha256:b1d1f0...5184d6",
        "linux/arm64": "sha256:fa8dc7...3d744b"
      }
    },
    "library/postgres": {
      "16-alpine": {
        "linux/amd64": "sha256:b0193a...4c27b1",
        "linux/arm64": "sha256:afa9bf...5e0e41"
      }
    }
  },
  "ghcr.io": {
    "home-assistant/home-assistant": {
      "latest": {
        "linux/amd64": "sha256:ef20dc...c940ca"
      }
    }
  }
}
```

### Nix Format

When using `--output-format=nix`, the application outputs a Nix attribute set that can be directly imported into Nix configurations:

```nix
{
  "docker.gitea.com" = {
    "gitea" = {
      "latest" = {
        "linux/amd64" = "sha256:5ee30f...de6367";
      };
    };
  };
  "docker.io" = {
    "library/busybox" = {
      "latest" = {
        "linux/amd64" = "sha256:ad9fa4...948f9f";
        "linux/arm/v7" = "sha256:b1d1f0...5184d6";
        "linux/arm64" = "sha256:fa8dc7...3d744b";
      };
    };
    "library/postgres" = {
      "16-alpine" = {
        "linux/amd64" = "sha256:b0193a...4c27b1";
        "linux/arm64" = "sha256:afa9bf...5e0e41";
      };
    };
  };
  "ghcr.io" = {
    "home-assistant/home-assistant" = {
      "latest" = {
        "linux/amd64" = "sha256:ef20dc...c940ca";
      };
    };
  };
}
```
