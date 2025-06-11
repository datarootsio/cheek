---
title: "Minimal Example"
description: "A simple example to get started with Cheek"
lead: "Get up and running with Cheek in minutes using this minimal configuration."
date: 2023-01-01T00:00:00+00:00
lastmod: 2023-01-01T00:00:00+00:00
draft: false
images: []
menu:
  docs:
    parent: "getting-started"
weight: 120
toc: true
---

## Quick Start

Here's a minimal example to get Cheek running:

### 1. Create a job configuration file

Create a file called `jobs.yaml`:

```yaml
jobs:
  hello:
    command: echo "Hello from Cheek!"
    cron: "*/5 * * * *"  # Run every 5 minutes
```

### 2. Run Cheek

```bash
cheek run jobs.yaml
```

That's it! Cheek will now execute the `hello` job every 5 minutes.

## What's happening?

- **jobs**: The root section that contains all your job definitions
- **hello**: The name of your job (you can choose any name)
- **command**: The shell command to execute
- **cron**: A cron expression that defines when to run the job

## Next Steps

- Add more jobs to your configuration
- Explore [job flow]({{< relref "job-flow" >}}) to chain jobs together
- Set up the [WebUI]({{< relref "webui" >}}) to monitor your jobs
- Learn about [events and notifications]({{< relref "events" >}})