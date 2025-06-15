---
title: Web UI
---

`cheek` ships with a web UI that by default gets launched on port `8081`. You can define the port on which it is accessible via the `--port` flag.

## Accessing the Web UI

You can access the UI by navigating to `http://localhost:8081`. When `cheek` is deployed you are recommended to NOT make this port publicly accessible, instead navigate to the UI via an SSH tunnel.

## Features

The UI allows you to:

- Get a quick overview of jobs that have run
- View jobs that have errored
- Access job logs and output
- Monitor job execution status in real-time

## Screenshots

![main-screen](/main.png)

![job-overview](/joboverview.png)

## Log Storage

The UI displays logs by fetching the state of the scheduler and by reading the logs that (per job) get written to the sqlite backend. Note that you can ignore these logs, as output of jobs will always go to stdout as well.

## Security Note

When `cheek` is deployed in production, you are recommended to NOT make the web UI port publicly accessible. Instead, access the UI via an SSH tunnel for security.