---
title: Installation
---

Fetch the latest version for your system below:

- [darwin-arm64](https://github.com/bart6114/cheek/releases/latest/download/cheek-darwin-arm64)
- [darwin-amd64](https://github.com/bart6114/cheek/releases/latest/download/cheek-darwin-amd64)
- [linux-386](https://github.com/bart6114/cheek/releases/latest/download/cheek-linux-386)
- [linux-arm64](https://github.com/bart6114/cheek/releases/latest/download/cheek-linux-arm64)
- [linux-amd64](https://github.com/bart6114/cheek/releases/latest/download/cheek-linux-amd64)

## Quick Installation

You can fetch and install `cheek` with the following commands:

```bash
curl -L https://github.com/bart6114/cheek/releases/latest/download/cheek-darwin-amd64 -o cheek
chmod +x cheek
./cheek
```

Optionally, put `cheek` on your `PATH` for system-wide access.

## Docker

Check out the `Dockerfile.example` for an example on how to use `cheek` within the context of a Docker container. Note that this builds upon a published Ubuntu-based image build that you can find in the base [Dockerfile](https://github.com/bart6114/cheek/blob/main/Dockerfile).

Prebuilt images are available at `ghcr.io/bart6114/cheek:latest` where `latest` can be replaced by a version tag. Check out the [available images](https://github.com/bart6114/cheek/pkgs/container/cheek) for an overview on available tags.

## Available Versions

If you want to pin your setup to a specific version of `cheek` you can use the following template to fetch your `cheek` binary:

- latest version: https://github.com/bart6114/cheek/releases/latest/download/cheek-{os}-{arch}
- tagged version: https://github.com/bart6114/cheek/releases/download/{tag}/cheek-{os}-{arch}

Where:

- `os` is one of `linux`, `darwin`
- `arch` is one of `amd64`, `arm64`, `386`
- `tag` is one the [available tags](https://github.com/bart6114/cheek/tags)