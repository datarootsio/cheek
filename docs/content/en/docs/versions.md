---
title: Available Versions
---

If you want to pin your setup to a specific version of `cheek` you can use the following template to fetch your `cheek` binary:

## Download URL Templates

- **Latest version**: `https://github.com/bart6114/cheek/releases/latest/download/cheek-{os}-{arch}`
- **Tagged version**: `https://github.com/bart6114/cheek/releases/download/{tag}/cheek-{os}-{arch}`

Where:

- `os` is one of `linux`, `darwin`
- `arch` is one of `amd64`, `arm64`, `386`
- `tag` is one the [available tags](https://github.com/bart6114/cheek/tags)

## Examples

### Latest Version Downloads
- Linux AMD64: `https://github.com/bart6114/cheek/releases/latest/download/cheek-linux-amd64`
- macOS ARM64: `https://github.com/bart6114/cheek/releases/latest/download/cheek-darwin-arm64`
- Linux ARM64: `https://github.com/bart6114/cheek/releases/latest/download/cheek-linux-arm64`

### Specific Version Downloads
Replace `{tag}` with your desired version (e.g., `v1.0.0`):
- `https://github.com/bart6114/cheek/releases/download/v1.0.0/cheek-linux-amd64`