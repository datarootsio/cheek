---
title: Overview
---

![cheek](/cheek-v2.svg)

`cheek` is a pico-sized declarative job scheduler designed to excel in a single-node environment. `cheek` aims to be lightweight, stand-alone and simple. It does not compete for robustness.

## What is cheek?

`cheek` is a lightweight job scheduler designed to run tasks on a schedule with minimal configuration. It provides a web interface for monitoring jobs and supports webhooks for notifications.

## Key Features

- **Lightweight**: Minimal resource footprint
- **Declarative**: Define your jobs and schedules in YAML
- **Single-node**: Optimized for single-node environments
- **Web UI**: Built-in web interface for monitoring and management
- **Event-driven**: Support for webhooks and job triggering
- **Simple**: Easy to set up and configure

## Why another scheduler?

For many simple cases, **crontab** would fit my bill for about 80% of scheduling needs. It's battle-tested, widely available, and gets the job done for basic periodic tasks.

However, I wanted something that was:

- **Declarative** - Define jobs and schedules in a clear, version-controlled format
- **Code-friendly** - Easy to manage through code updates and deployments
- **Log accessible** - More intuitive access to job logs and execution history
- **Lightweight** - Simple to deploy and maintain, without enterprise complexity

Cheek emerged from this need for a **very lightweight scheduler with slightly more batteries included than crontab**. It bridges the gap between crontab's simplicity and full-featured enterprise schedulers.

Think of Cheek as crontab's slightly more sophisticated cousin - retaining the simplicity while adding just enough modern conveniences to make job management more pleasant.

