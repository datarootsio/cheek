---
title: Docker
---

## Docker Usage

Check out the `Dockerfile.example` for an example on how to use `cheek` within the context of a Docker container. Note that this builds upon a published Ubuntu-based image build that you can find in the base [Dockerfile](https://github.com/bart6114/cheek/blob/main/Dockerfile).

## Prebuilt Images

Prebuilt images are available at `ghcr.io/bart6114/cheek:latest` where `latest` can be replaced by a version tag. 

Check out the [available images](https://github.com/bart6114/cheek/pkgs/container/cheek) for an overview on available tags.

## Running with Docker

You can run cheek using Docker with your configuration file:

```bash
docker run -v $(pwd)/my-schedule.yaml:/schedule.yaml ghcr.io/bart6114/cheek:latest run /schedule.yaml
```

Make sure to mount your schedule configuration file into the container so cheek can access it.