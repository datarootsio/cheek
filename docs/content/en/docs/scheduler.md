---
title: Scheduler
---

The core of `cheek` consists of a scheduler that uses the schedule specs defined in your `yaml` file to trigger jobs when they are due.

## Running the Scheduler

You can launch the scheduler via:

```bash
cheek run ./path/to/my-schedule.yaml
```

Check out `cheek run --help` for configuration options.

## How It Works

The scheduler continuously monitors your job definitions and executes them according to their cron schedules. Jobs are executed in separate processes, allowing for concurrent execution unless specifically disabled.

## Job Features

- **Cron Scheduling**: Use standard cron expressions to define when jobs run
- **Retries**: Configure automatic retries for failed jobs
- **Concurrent Execution Control**: Prevent multiple instances of the same job from running simultaneously
- **Working Directory**: Specify custom working directories for jobs
- **Environment Variables**: Set custom environment variables for each job
- **Job Triggering**: Trigger other jobs based on success or failure events