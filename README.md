<img src="docs/static/cheek-v2.svg" alt="cheek" />

# cheek

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/bart6114/cheek?label=version)
![workflow](https://github.com/bart6114/cheek/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bart6114/cheek)](https://goreportcard.com/report/github.com/bart6114/cheek)
[![Go Reference](https://pkg.go.dev/badge/github.com/bart6114/cheek.svg)](https://pkg.go.dev/github.com/bart6114/cheek)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
![love](https://img.shields.io/badge/made_with-%E2%9D%A4%EF%B8%8F-blue)

`cheek` is a pico-sized declarative job scheduler designed to excel in a single-node environment. `cheek` aims to be lightweight, stand-alone and simple. It does not compete for robustness.

## Quick Start

Fetch the latest version for your system:

[darwin-arm64](https://github.com/bart6114/cheek/releases/latest/download/cheek-darwin-arm64) |
[darwin-amd64](https://github.com/bart6114/cheek/releases/latest/download/cheek-darwin-amd64) |
[linux-386](https://github.com/bart6114/cheek/releases/latest/download/cheek-linux-386) |
[linux-arm64](https://github.com/bart6114/cheek/releases/latest/download/cheek-linux-arm64) |
[linux-amd64](https://github.com/bart6114/cheek/releases/latest/download/cheek-linux-amd64)

```sh
curl -L https://github.com/bart6114/cheek/releases/latest/download/cheek-darwin-amd64 -o cheek
chmod +x cheek
./cheek run schedule.yaml
```

## Documentation

For comprehensive documentation, configuration options, examples, and advanced features, visit:

**📚 [cheek.barts.space](https://cheek.barts.space)**

## Acknowledgements

`cheek` is building on top of many great OSS assets. Notable thanks goes to:

- [gronx](https://github.com/adhocore/gronx): for allowing me not to worry about CRON strings.

<br/>
 
![GitHub Contributors](https://contrib.rocks/image?repo=bart6114/cheek)
